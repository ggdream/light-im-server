package conversation

import (
	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/pkg/cache"
	"lim/pkg/errno"
)

type deleteReqModel struct {
	ConversationID *string `json:"conversation_id" binding:"required"`
}

func DeleteController(c *gin.Context) {
	var (
		form deleteReqModel
		err  error
	)
	if err = c.ShouldBindJSON(&form); err != nil {
		errno.NewFParamInvalid(err.Error()).Reply(c)
		return
	}

	userId := config.CtxKeyManager.GetUserID(c)
	ca := cache.ChatConv{}
	err = ca.Del(c.Request.Context(), userId, *form.ConversationID)
	if err != nil {
		errno.NewF(errno.BaseErrRedis, err.Error(), errno.ErrChatConvDelFailed).Reply(c)
		return
	}

	errno.NewS(nil).Reply(c)
}
