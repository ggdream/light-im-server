package conversation

import (
	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/pkg/cache"
	"lim/pkg/errno"
)

type listResModel struct {
	Items []listResModelItem `json:"items"`
}

type listResModelItem struct {
	SenderID       string `json:"sender_id"`
	ReceiverID     string `json:"receiver_id"`
	ConversationID string `json:"conversation_id"`
	Type           uint8  `json:"type"`
	Text           string `json:"text"`
	Image          string `json:"image"`
	Audio          string `json:"audio"`
	Video          string `json:"video"`
	Custom         string `json:"custom"`
	IsRead         uint8  `json:"is_read"`
	Unread         int64  `json:"unread"`
	Timestamp      int64  `json:"timestamp"`
	Sequence       int64  `json:"sequence"`
	CreateAt       int64  `json:"create_at"`
}

func ListContorller(c *gin.Context) {
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

	items := make([]listResModelItem, 0, len(res))
	for _, v := range res {
		items = append(items, listResModelItem{
			SenderID:       v.SenderID,
			ReceiverID:     v.ReceiverID,
			ConversationID: v.ConversationID,
			Type:           v.Type,
			Text:           v.Text,
			Image:          v.Image,
			Audio:          v.Audio,
			Video:          v.Video,
			Custom:         v.Custom,
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

	ret := &listResModel{
		Items: items,
	}
	errno.NewS(ret).Reply(c)
}
