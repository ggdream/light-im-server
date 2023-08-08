package message

import (
	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/pkg/cache"
	"lim/pkg/db"
	"lim/pkg/errno"
)

type markReqModel struct {
	UserID   *string `json:"user_id" binding:"required"`
	Sequence *int64  `json:"sequence" binding:"required"`
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
	err = doc.MarkAsRead(*form.UserID, userId, *form.Sequence)
	if err != nil {
		errno.NewF(errno.BaseErrMongo, err.Error(), errno.ErrAuthLoginFailed).Reply(c)
		return
	}

	ca := cache.ChatConv{}
	_ = ca.MarkAsRead(userId, *form.UserID)

	// TODO: hub通知已读

	errno.NewS(nil).Reply(c)
}
