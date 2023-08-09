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
	clientGroup := rg.Group("/c")
	{
		connectGroup := clientGroup.Group("/connect")
		connectGroup.GET("/im", socket.IM)
		connectGroup.POST("/logout", auth.LogoutController)

		messageGroup := clientGroup.Group("/message")
		messageGroup.POST("/pull", message.PullController)     // c
		messageGroup.POST("/mark", message.MarkController)        // c
		messageGroup.POST("/send", message.SendController(false)) // c s

		convGroup := clientGroup.Group("/conv")
		convGroup.POST("/pull", conversation.PullContorller)     // c
		convGroup.POST("/delete", conversation.DeleteController) // c
	}

	serverGroup := rg.Group("/s")
	{
		authGroup := serverGroup.Group("/auth")
		authGroup.POST("/sign", auth.SignController) // s

		userGroup := rg.Group("/user")
		userGroup.POST("/create", user.CreateController) // s
		userGroup.POST("/delete", user.DeleteController) // s
		userGroup.POST("/update", user.UpdateController) // s

		messageGroup := rg.Group("/message")
		messageGroup.POST("/send", message.SendController(true)) // c s
	}
}
