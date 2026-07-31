package router

import (
	"clover_server/api"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
)

func GlobalNotificationRouter(rg *gin.RouterGroup) {
	app := api.App.GlobalNotificationApi
	rg.POST("global_notification", middleware.AdminAuthMiddleware, app.CreateView)
	rg.GET("global_notification", middleware.AuthMiddleware, app.ListView)
	rg.DELETE("global_notification", middleware.AdminAuthMiddleware, app.RemoveAdminView)
	rg.POST("global_notification/user", middleware.AuthMiddleware, app.UserMsgActionView)
}
