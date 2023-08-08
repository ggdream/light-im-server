package auth

import (
	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/pkg/cache"
	"lim/pkg/errno"
)

func LogoutController(c *gin.Context) {
	var (
		ca  cache.AuthToken
		err error
	)

	userId := config.CtxKeyManager.GetUserID(c)
	err = ca.Del(userId, 0)
	if err != nil {
		errno.NewF(errno.BaseErrRedis, err.Error(), errno.ErrAuthLogoutFailed).Reply(c)
		return
	}

	errno.NewS(nil).Reply(c)
}
