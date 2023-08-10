package middleware

import (
	"unsafe"

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
		if exist.ExistInSlice(c.Request.URL.Path, []string{"/auth/sign", "/im"}) {
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

		tokenKey := unsafe.Slice(unsafe.StringData(config.GetApp().TokenKey), len(config.GetApp().TokenKey))
		userId, err := token.Parse(tokenKey, header.Token)
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

		isPass, err := authTokenCa.Verify(userId, header.Token)
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
