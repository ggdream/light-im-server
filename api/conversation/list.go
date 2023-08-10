package conversation

import (
	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/pkg/cache"
	"lim/pkg/errno"
)

type pullResModel struct {
	Count int                `json:"count"`
	Items []pullResModelItem `json:"items"`
}

type pullResModelItem struct {
	SenderID       string `json:"sender_id"`
	ReceiverID     string `json:"receiver_id"`
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	Type           uint8  `json:"type"`
	Text           string `json:"text"`
	Image          string `json:"image"`
	Audio          string `json:"audio"`
	Video          string `json:"video"`
	Custom         string `json:"custom"`
	IsSelf         uint8  `json:"is_self"`
	IsRead         uint8  `json:"is_read"`
	Unread         int64  `json:"unread"`
	Timestamp      int64  `json:"timestamp"`
	Sequence       int64  `json:"sequence"`
	CreateAt       int64  `json:"create_at"`
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
			Text:           v.Text,
			Image:          v.Image,
			Audio:          v.Audio,
			Video:          v.Video,
			Custom:         v.Custom,
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
