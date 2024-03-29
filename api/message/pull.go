package message

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"lim/config"
	"lim/pkg/db"
	"lim/pkg/errno"
)

type pullReqModel struct {
	ConversationID *string `json:"conversation_id" binding:"required"`
	Sequence       int64   `json:"sequence"`
}

type pullResModel struct {
	Sequence int64              `json:"sequence"`
	IsEnd    uint8              `json:"is_end"`
	Items    []pullResModelItem `json:"items"`
}

type pullResModelItem struct {
	SenderID       string         `json:"sender_id"`
	ReceiverID     string         `json:"receiver_id"`
	UserID         string         `json:"user_id"`
	GroupID        string         `json:"group_id"`
	ConversationID string         `json:"conversation_id"`
	Type           uint8          `json:"type"`
	Text           *db.TextElem   `json:"text"`
	Image          *db.ImageElem  `json:"image"`
	Audio          *db.AudioElem  `json:"audio"`
	Video          *db.VideoElem  `json:"video"`
	Record         *db.RecordElem `json:"record"`
	Custom         *db.CustomElem `json:"custom"`
	IsSelf         uint8          `json:"is_self"`
	IsRead         uint8          `json:"is_read"`
	IsPeerRead     uint8          `json:"is_peer_read"`
	Timestamp      int64          `json:"timestamp"`
	Sequence       int64          `json:"sequence"`
	CreateAt       int64          `json:"create_at" copier:"Cts"`
}

func PullController(c *gin.Context) {
	var (
		form pullReqModel
		err  error
	)
	if err = c.ShouldBindJSON(&form); err != nil {
		errno.NewFParamInvalid(err.Error()).Reply(c)
		return
	}

	userId := config.CtxKeyManager.GetUserID(c)
	doc := db.ChatMsgDoc{}
	res, err := doc.List(c, *form.ConversationID, userId, form.Sequence, 20)
	if err != nil {
		errno.NewF(errno.BaseErrMongo, err.Error(), errno.ErrAuthLoginFailed).Reply(c)
		return
	}

	items := make([]pullResModelItem, 0, len(res))
	err = copier.Copy(&items, &res)
	if err != nil {
		errno.NewF(errno.BaseErrTools, err.Error(), errno.ErrCopier).Reply(c)
		return
	}

	for i, v := range items {
		if v.SenderID == userId {
			items[i].UserID = v.ReceiverID
			items[i].IsPeerRead = items[i].IsRead
			items[i].IsRead = 1
			items[i].IsSelf = 1
		} else {
			items[i].UserID = v.SenderID
			items[i].IsPeerRead = 1
			items[i].IsSelf = 0
		}
	}

	var (
		sequence int64 = 0
		isEnd    uint8 = 0
	)
	if len(items) != 0 {
		sequence = items[len(items)-1].Sequence
	}
	if len(items) < 20 {
		isEnd = 1
	}

	ret := &pullResModel{
		Sequence: sequence,
		IsEnd:    isEnd,
		Items:    items,
	}
	errno.NewS(ret).Reply(c)
}
