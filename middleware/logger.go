package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"lim/config"
	"lim/tools/log"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		e := config.CtxKeyManager.GetError(c)
		model := log.BusinessModel{
			Type:   log.LogTypeMiddleWare,
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Time:   startTime.Format(time.DateTime),
			Cost:   time.Now().Sub(startTime).Milliseconds(),
			UserID: config.CtxKeyManager.GetUserID(c),
			UserIP: c.ClientIP(),
			// Role:   c.GetUint(global.CtxKeyUserRole),
			Error:  e,
		}
		log.InfoWf("", model.Encode())
	}
}
