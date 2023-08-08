package socket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"lim/config"
	"lim/hub"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  0,
		WriteBufferSize: 0,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

func IM(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	userId := config.CtxKeyManager.GetUserID(c)
	hub.SetConn2Hub(userId, conn)
}
