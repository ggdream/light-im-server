package message

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"lim/config"
	"lim/pkg/db"
	"lim/pkg/errno"
)

type historyReqModel struct {
	UserID   *string `json:"user_id" binding:"required"`
	Sequence int64   `json:"sequence"`
}

type historyResModel struct {
	Sequence int64                 `json:"sequence"`
	IsEnd    uint8                 `json:"is_end"`
	Items    []historyResModelItem `json:"items"`
}

type historyResModelItem struct {
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
	CreateAt       int64  `json:"create_at" bson:"create_at" copier:"cts"`
}

func HistoryController(c *gin.Context) {
	var (
		form historyReqModel
		err  error
	)
	if err = c.ShouldBindJSON(&form); err != nil {
		errno.NewFParamInvalid(err.Error()).Reply(c)
		return
	}

	userId := config.CtxKeyManager.GetUserID(c)
	doc := db.ChatMsgDoc{}
	res, err := doc.List(*form.UserID, userId, form.Sequence, 20)
	if err != nil {
		errno.NewF(errno.BaseErrMongo, err.Error(), errno.ErrAuthLoginFailed).Reply(c)
		return
	}

	items := make([]historyResModelItem, 0, len(res))
	err = copier.Copy(&items, &res)
	if err != nil {
		errno.NewF(errno.BaseErrTools, err.Error(), errno.ErrCopier).Reply(c)
		return
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

	ret := &historyResModel{
		Sequence: sequence,
		IsEnd:    isEnd,
		Items:    items,
	}
	errno.NewS(ret).Reply(c)
}
