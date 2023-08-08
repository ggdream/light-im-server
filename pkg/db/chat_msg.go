package db

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatMsgDoc struct {
	SenderID       string `json:"sender_id" bson:"sender_id"`
	ReceiverID     string `json:"receiver_id" bson:"receiver_id"`
	ConversationID string `json:"conversation_id" bson:"conversation_id"`
	Type           uint8  `json:"type" bson:"type"`
	Text           string `json:"text" bson:"text"`
	Image          string `json:"image" bson:"image"`
	Audio          string `json:"audio" bson:"audio"`
	Video          string `json:"video" bson:"video"`
	Custom         string `json:"custom" bson:"custom"`
	IsRead         uint8  `json:"is_read" bson:"is_read"`
	Timestamp      int64  `json:"timestamp" bson:"timestamp"`
	Sequence       int64  `json:"sequence" bson:"sequence"`
	CreateTs       int64  `json:"create_ts" bson:"create_ts" copier:"cts"`
	CreateAt       string `json:"create_at" bson:"create_at"`
	DeleteAt       string `json:"delete_at" bson:"delete_at"`
}

func (m *ChatMsgDoc) Create() error {
	m.ConversationID = m.genCID(m.SenderID, m.ReceiverID)
	m.IsRead = 0
	t := time.Now()
	m.CreateTs = t.UnixMilli()
	m.CreateAt = t.Format(time.RFC3339)

	return client.Insert(m.DocName(), m)
}

func (m *ChatMsgDoc) MarkAsRead(senderId, receiverId string, sequence int64) error {
	return client.UpdateMany(
		m.DocName(),
		bson.D{
			{
				Key: "sender_id", Value: senderId,
			},
			{
				Key: "receiver_id", Value: receiverId,
			},
			{
				Key: "sequence", Value: bson.M{
					"$lte": sequence,
				},
			},
		},
		bson.D{
			{
				Key: "is_read", Value: 1,
			},
		},
	)
}

func (m *ChatMsgDoc) List(senderId, receiverId string, sequence, number int64) ([]ChatMsgDoc, error) {
	option := options.Find().SetSort(map[string]interface{}{"_id": -1}).SetLimit(number)
	filter := bson.D{
		{
			Key: "conversation_id", Value: m.genCID(senderId, receiverId),
		},
	}
	if sequence != 0 {
		filter = append(filter, bson.E{
			Key: "sequence", Value: bson.M{
				"$lt": sequence,
			},
		})
	}

	cursor, err := client.SearchMany(m.DocName(), filter, option)
	if err != nil {
		return nil, err
	}

	var res []ChatMsgDoc
	err = cursor.All(nil, &res)
	if err != nil {
		return nil, err
	}
	
	for i, v := range res {
		if v.SenderID == receiverId {
			res[i].IsRead = 1
		}
	}

	return res, nil
}

func (m *ChatMsgDoc) genCID(a, b string) string {
	if strings.Compare(a, b) == 1 {
		a, b = b, a
	}

	return fmt.Sprintf("c_%s_%s", a, b)
}

func (m *ChatMsgDoc) DocName() string { return "lim_chat_msg" }
