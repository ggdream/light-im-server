package middleware

import (
	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/pkg/cache"
	"lim/pkg/errno"
	"lim/tools/exist"
	"lim/tools/token"
)

type authTokenHeader struct {
	Token string `header:"light-im-token" binding:"required"`
}

func AuthToken() gin.HandlerFunc {
	var (
		authTokenCa cache.AuthToken
	)

	return func(c *gin.Context) {
		if exist.ExistInSlice(c.Request.URL.Path, []string{"/auth/login"}) {
			c.Next()

			return
		}

		var (
			header authTokenHeader
			err    error
		)
		if err = c.ShouldBindHeader(&header); err != nil {
			errno.Fail(c, errno.ErrJWT, err.Error())
			return
		}

		userId, role, err := token.Parse(config.GetApp().TokenKey, header.Token)
		if err != nil {
			switch err {
			case token.ErrTokenExpired:
				errno.NewF(errno.BaseErrInvalid, err.Error(), errno.ErrJWT).Reply(c)
			case token.ErrTokenMalformed:
				errno.NewF(errno.BaseErrInvalid, err.Error(), errno.ErrJWT).Reply(c)
			default:
				errno.NewF(errno.BaseErrInvalid, err.Error(), errno.ErrJWT).Reply(c)
			}
			return
		}

		isPass, err := authTokenCa.Verify(userId, role, header.Token)
		if err != nil {
			errno.NewF(errno.BaseErrRedis, err.Error(), errno.ErrRedis).Reply(c)
			return
		}
		if !isPass {
			errno.NewF(errno.BaseErrInvalid, "登录凭证失效", errno.ErrJWT).Reply(c)
			return
		}


		config.CtxKeyManager.SetUserID(c, userId)
		c.Next()
	}
}
