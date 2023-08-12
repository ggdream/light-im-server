package conversation

import (
	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/pkg/cache"
	"lim/pkg/db"
	"lim/pkg/errno"
	"lim/tools/typec"
)

type pullResModel struct {
	Count int                `json:"count"`
	Items []pullResModelItem `json:"items"`
}

type pullResModelItem struct {
	SenderID       string         `json:"sender_id"`
	ReceiverID     string         `json:"receiver_id"`
	UserID         string         `json:"user_id"`
	ConversationID string         `json:"conversation_id"`
	Type           uint8          `json:"type"`
	Text           *db.TextElem   `json:"text"`
	Image          *db.ImageElem  `json:"image"`
	Audio          *db.AudioElem  `json:"audio"`
	Video          *db.VideoElem  `json:"video"`
	File           *db.FileElem   `json:"file"`
	Custom         *db.CustomElem `json:"custom"`
	Record         *db.RecordElem `json:"record"`
	IsSelf         uint8          `json:"is_self"`
	IsRead         uint8          `json:"is_read"`
	Unread         int64          `json:"unread"`
	Timestamp      int64          `json:"timestamp"`
	Sequence       int64          `json:"sequence"`
	CreateAt       int64          `json:"create_at"`
}

func PullContorller(c *gin.Context) {
	var (
		ca  cache.ChatConv
		err error
	)

	userId := config.CtxKeyManager.GetUserID(c)
	res, err := ca.List(userId)
	if err != nil {
		errno.NewF(errno.BaseErrRedis, err.Error(), errno.ErrChatConvListFailed).Reply(c)
		return
	}

	items := make([]pullResModelItem, 0, len(res))
	for _, v := range res {
		var isSelf uint8 = 0
		uid := v.SenderID
		if v.SenderID == userId {
			isSelf = 1
			uid = v.ReceiverID
		}

		items = append(items, pullResModelItem{
			SenderID:       v.SenderID,
			ReceiverID:     v.ReceiverID,
			UserID:         uid,
			ConversationID: v.ConversationID,
			Type:           v.Type,
			Text:           typec.JsonToStruct[db.TextElem](v.Text),
			Image:          typec.JsonToStruct[db.ImageElem](v.Image),
			Audio:          typec.JsonToStruct[db.AudioElem](v.Audio),
			Video:          typec.JsonToStruct[db.VideoElem](v.Video),
			File:           typec.JsonToStruct[db.FileElem](v.File),
			Custom:         typec.JsonToStruct[db.CustomElem](v.Custom),
			Record:         typec.JsonToStruct[db.RecordElem](v.Record),
			IsSelf:         isSelf,
			IsRead: func() uint8 {
				if v.Unread == 0 {
					return 0
				}

				return 1
			}(),
			Unread:    v.Unread,
			Timestamp: v.Timestamp,
			Sequence:  v.Sequence,
			CreateAt:  v.CreateAt,
		})
	}

	ret := &pullResModel{
		Count: len(items),
		Items: items,
	}
	errno.NewS(ret).Reply(c)
}
