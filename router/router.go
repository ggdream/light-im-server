package router

import (
	"github.com/gin-gonic/gin"

	"lim/api/auth"
	"lim/api/conversation"
	"lim/api/message"
	"lim/api/socket"
	"lim/api/user"
)

func Set(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	authGroup.POST("/login", auth.LoginController)
	authGroup.POST("/logout", auth.LogoutController)

	convGroup := rg.Group("/conv")
	convGroup.POST("/list", conversation.ListContorller)
	convGroup.POST("/delete", conversation.DeleteController)

	messageGroup := rg.Group("/message")
	messageGroup.POST("/list", message.HistoryController)
	messageGroup.POST("/mark", message.MarkController)
	messageGroup.POST("/send", message.SendController())

	rg.GET("/im", socket.IM)

	userGroup := rg.Group("/user")
	userGroup.POST("/create", user.CreateController)
	userGroup.POST("/delete", user.DeleteController)
	userGroup.POST("/update", user.UpdateController)
}
