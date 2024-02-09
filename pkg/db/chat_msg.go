package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	clone "github.com/huandu/go-clone/generic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"lim/config"
)

type ChatMsgDoc struct {
	BaseAutoIDModel `bson:",inline"`
	UserID          string      `json:"user_id" bson:"user_id"`
	SenderID        string      `json:"sender_id" bson:"sender_id"`
	ReceiverID      string      `json:"receiver_id" bson:"receiver_id"`
	GroupID         string      `json:"group_id" bson:"group_id"`
	ConversationID  string      `json:"conversation_id" bson:"conversation_id"`
	Type            uint8       `json:"type" bson:"type"`
	Text            *TextElem   `json:"text" bson:"text"`
	Image           *ImageElem  `json:"image" bson:"image"`
	Audio           *AudioElem  `json:"audio" bson:"audio"`
	Video           *VideoElem  `json:"video" bson:"video"`
	File            *FileElem   `json:"file" bson:"file"`
	Custom          *CustomElem `json:"custom" bson:"custom"`
	Record          *RecordElem `json:"record" bson:"record"`
	IsRead          uint8       `json:"is_read" bson:"is_read"`
	MarkAt          int64       `json:"mark_at" bson:"mark_at"`
	Timestamp       int64       `json:"timestamp" bson:"timestamp"`
	Sequence        int64       `json:"sequence" bson:"sequence"`
}

type TextElem struct {
	Text string `json:"text" bson:"text"`
}

type ImageElem struct {
	Name         string `json:"name" bson:"name"`
	Size         int64  `json:"size" bson:"size"`
	ContentType  string `json:"content_type" bson:"content_type"`
	URL          string `json:"url" bson:"url"`
	ThumbnailURL string `json:"thumbnail_url" bson:"thumbnail_url"`
}

type AudioElem struct {
	Name        string `json:"name" bson:"name"`
	Size        int64  `json:"size" bson:"size"`
	ContentType string `json:"content_type" bson:"content_type"`
	Duration    int64  `json:"duration" bson:"duration"`
	URL         string `json:"url" bson:"url"`
}

type VideoElem struct {
	Name         string `json:"name" bson:"name"`
	Size         int64  `json:"size" bson:"size"`
	ContentType  string `json:"content_type" bson:"content_type"`
	Duration     int64  `json:"duration" bson:"duration"`
	URL          string `json:"url" bson:"url"`
	ThumbnailURL string `json:"thumbnail_url" bson:"thumbnail_url"`
}

type FileElem struct {
	Name        string `json:"name" bson:"name"`
	Size        int64  `json:"size" bson:"size"`
	ContentType string `json:"content_type" bson:"content_type"`
	URL         string `json:"url" bson:"url"`
}

type CustomElem struct {
	Content string `json:"content" bson:"content"`
}

type RecordElem struct {
	Size        int64  `json:"size" bson:"size"`
	ContentType string `json:"content_type" bson:"content_type"`
	Duration    int64  `json:"duration" bson:"duration"`
	URL         string `json:"url" bson:"url"`
}

func (m *ChatMsgDoc) IsPrivate() bool {
	return m.ReceiverID != "" || strings.HasPrefix(m.ConversationID, "c_")
}

func (m *ChatMsgDoc) IsGroup() bool {
	return m.GroupID != "" || strings.HasPrefix(m.ConversationID, "g_")
}

func (m *ChatMsgDoc) Create(ctx context.Context) ([]string, error) {
	m.ConversationID = m.GenCID()
	m.IsRead = 1
	m.CreateAt = time.Now().UnixMilli()
	m.Status = config.RecordStatusNormal

	var receiverIds []string

	_, err := client.Database().Collection(m.DocName()).InsertOne(ctx, m)
	if err != nil {
		return nil, err
	}

	// 单聊
	if m.IsPrivate() {
		m1 := clone.Slowly(m)
		m1.UserID = m.ReceiverID
		m1.IsRead = 0
		_, err = client.Database().Collection(m.DocName()).InsertOne(ctx, m1)
		if err != nil {
			return nil, err
		}

		receiverIds = append(receiverIds, m1.UserID)

		return receiverIds, nil
	}

	// 群聊
	cur, err := client.Database().Collection(GetGroupMemberDocName()).Find(ctx, bson.M{"_id": m.GroupID})
	if err != nil {
		return nil, err
	}
	var groupRecords []GroupMemberDoc
	err = cur.All(ctx, &groupRecords)
	if err != nil {
		return nil, err
	}

	var writeModels []mongo.WriteModel
	for _, record := range groupRecords {
		if record.UserID != m.UserID {
			m1 := clone.Slowly(m)
			m1.UserID = record.UserID
			m1.IsRead = 0
			writeModels = append(writeModels, mongo.NewInsertOneModel().SetDocument(m))
			receiverIds = append(receiverIds, record.UserID)
		}
	}
	_, err = client.Database().Collection(m.DocName()).BulkWrite(ctx, writeModels)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return receiverIds, nil
}

// func (m *ChatMsgDoc) Create(ctx context.Context) ([]string, error) {
// 	m.ID = strconv.Itoa(int(m.Sequence))
// 	m.ConversationID = m.GenCID()
// 	m.IsRead = 1
// 	m.CreateAt = time.Now().UnixMilli()
// 	m.Status = config.RecordStatusNormal

// 	var receiverIds []string
// 	err := client.Client().UseSession(ctx, func(sCtx mongo.SessionContext) (err error) {
// 		err = sCtx.StartTransaction()
// 		if err != nil {
// 			return
// 		}

// 		defer func() {
// 			if err != nil {
// 				_ = sCtx.AbortTransaction(ctx)
// 			}
// 		}()

// 		_, err = client.Database().Collection(m.DocName()).InsertOne(sCtx, m)
// 		if err != nil {
// 			return err
// 		}

// 		// 单聊
// 		if m.IsPrivate() {
// 			m1 := clone.Slowly(m)
// 			m1.UserID = m.ReceiverID
// 			m1.IsRead = 0
// 			_, err = client.Database().Collection(m.DocName()).InsertOne(sCtx, m1)
// 			if err != nil {
// 				return err
// 			}

// 			receiverIds = append(receiverIds, m1.UserID)

// 			return sCtx.CommitTransaction(ctx)
// 		}

// 		// 群聊
// 		cur, err := client.Database().Collection(GetGroupMemberDocName()).Find(sCtx, bson.M{"_id": m.GroupID})
// 		if err != nil {
// 			return err
// 		}
// 		var groupRecords []GroupMemberDoc
// 		err = cur.All(sCtx, &groupRecords)
// 		if err != nil {
// 			return err
// 		}

// 		var writeModels []mongo.WriteModel
// 		for _, record := range groupRecords {
// 			if record.UserID != m.UserID {
// 				m1 := clone.Slowly(m)
// 				m1.UserID = record.UserID
// 				m1.IsRead = 0
// 				writeModels = append(writeModels, mongo.NewInsertOneModel().SetDocument(m))
// 				receiverIds = append(receiverIds, record.UserID)
// 			}
// 		}
// 		_, err = client.Database().Collection(m.DocName()).BulkWrite(sCtx, writeModels)
// 		if err != nil {
// 			return err
// 		}

// 		return sCtx.CommitTransaction(ctx)
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return receiverIds, nil
// }

func (m *ChatMsgDoc) MarkAsRead(ctx context.Context, conversationId, userId string, sequence int64) error {
	_, err := client.Database().Collection(m.DocName()).UpdateMany(
		ctx,
		bson.D{
			{
				Key: "user_id", Value: userId,
			},
			{
				Key: "conversation_id", Value: conversationId,
			},
			{
				Key: "sequence", Value: bson.M{
					"$lte": sequence,
				},
			},
			{
				Key: "is_read", Value: 0,
			},
		},
		bson.D{
			{
				Key: "$set", Value: bson.M{
					"is_read": 1,
					"mark_at": time.Now().UnixMilli(),
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *ChatMsgDoc) List(ctx context.Context, conversationId, userId string, sequence, number int64) ([]ChatMsgDoc, error) {
	option := options.Find().SetSort(map[string]interface{}{"_id": -1}).SetLimit(number)
	filter := bson.D{
		{
			Key: "user_id", Value: userId,
		},
		{
			Key: "conversation_id", Value: conversationId,
		},
	}
	if sequence != 0 {
		filter = append(filter, bson.E{
			Key: "sequence", Value: bson.M{
				"$lt": sequence,
			},
		})
	}

	cursor, err := client.Database().Collection(m.DocName()).Find(ctx, filter, option)
	if err != nil {
		return nil, err
	}

	var res []ChatMsgDoc
	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *ChatMsgDoc) CountUnread(ctx context.Context, userId string) (int64, error) {
	return client.Database().Collection(m.DocName()).CountDocuments(ctx, bson.M{
		"user_id": userId,
		"is_read": 0,
	})
}

func (m *ChatMsgDoc) GenCID() string {
	if m.GroupID == "" {
		a, b := m.SenderID, m.ReceiverID
		if strings.Compare(a, b) == 1 {
			a, b = b, a
		}

		return fmt.Sprintf("c_%s_%s", a, b)
	}

	return fmt.Sprintf("g_%s", m.GroupID)
}

func (m *ChatMsgDoc) DocName() string { return "lim_chat_msg" }
