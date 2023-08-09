package auth

import (
	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/hub"
	"lim/pkg/cache"
	"lim/pkg/errno"
)

func LogoutController(c *gin.Context) {
	var (
		ca  cache.AuthToken
		err error
	)

	userId := config.CtxKeyManager.GetUserID(c)
	err = ca.Del(userId)
	if err != nil {
		errno.NewF(errno.BaseErrRedis, err.Error(), errno.ErrAuthLogoutFailed).Reply(c)
		return
	}

	hub.DelConn4Hub(userId)

	errno.NewS(nil).Reply(c)
}
