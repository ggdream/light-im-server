package message

import (
	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/pkg/db"
	"lim/pkg/errno"
)

type unreadResModel struct {
	Count int64 `json:"count"`
}

func UnreadController(c *gin.Context) {
	var (
		doc db.ChatMsgDoc
		err error
	)

	userId := config.CtxKeyManager.GetUserID(c)
	res, err := doc.CountUnrad(userId)
	if err != nil {
		errno.NewF(errno.BaseErrMongo, err.Error(), errno.ErrChatMsgUnreadFailed).Reply(c)
		return
	}

	ret := &unreadResModel{
		Count: res,
	}
	errno.NewS(ret).Reply(c)
}
