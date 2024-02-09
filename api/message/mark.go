package message

import (
	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/pkg/cache"
	"lim/pkg/db"
	"lim/pkg/errno"
)

type markReqModel struct {
	ConversationID *string `json:"conversation_id" binding:"required"`
	Sequence       *int64  `json:"sequence" binding:"required"`
}

func MarkController(c *gin.Context) {
	var (
		form markReqModel
		err  error
	)
	if err = c.ShouldBindJSON(&form); err != nil {
		errno.NewFParamInvalid(err.Error()).Reply(c)
		return
	}

	userId := config.CtxKeyManager.GetUserID(c)
	doc := db.ChatMsgDoc{}
	err = doc.MarkAsRead(c, *form.ConversationID, userId, *form.Sequence)
	if err != nil {
		errno.NewF(errno.BaseErrMongo, err.Error(), errno.ErrChatMsgMarkFailed).Reply(c)
		return
	}

	ca := cache.ChatConv{}
	_ = ca.MarkAsRead(c.Request.Context(), userId, *form.ConversationID)

	// TODO: hub通知已读

	errno.NewS(nil).Reply(c)
}
