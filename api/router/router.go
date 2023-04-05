package router

import (
	"github.com/gin-gonic/gin"
	"mygochat/api/handler"
)

func Register() *gin.Engine {
	r := gin.Default()
	initUserRouter(r)
	initChatRouter(r)
	r.NoRoute()
	return r
}

func initUserRouter(r *gin.Engine) {
	userGroup := r.Group("/user")
	userGroup.POST("/register", handler.Register)
	userGroup.POST("/login", handler.Login)
}

func initChatRouter(r *gin.Engine) {
	chatGroup := r.Group("/chat")
	chatGroup.Use(handler.CheckAuth)
	{
		chatGroup.POST("/person", handler.Chat)
		chatGroup.POST("/room", handler.ChatRoom)
	}
}
