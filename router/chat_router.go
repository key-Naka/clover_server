package router

import (
	"clover_server/api"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
)

func ChatRouter(rg *gin.RouterGroup) {
	app := api.App.ChatApi
	rg.GET("chat", middleware.AuthMiddleware, app.ChatListView)
	rg.GET("chat/session", middleware.AuthMiddleware, app.SessionListView)
	rg.DELETE("chat", middleware.AuthMiddleware, app.UserChatDeleteView)
	rg.DELETE("chat/user/:id", middleware.AuthMiddleware, app.UserChatDeleteByUserView)
	rg.POST("chat/read/:id", middleware.AuthMiddleware, app.ChatReadView)
	rg.GET("chat/ws", app.ChatView)
}
