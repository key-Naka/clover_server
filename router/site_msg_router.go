package router

import (
	"clover_server/api"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
)

func SiteMsgRouter(rg *gin.RouterGroup) {
	app := api.App.SiteMsgApi
	rg.GET("site_msg", middleware.AuthMiddleware, app.SiteMsgListView)
	rg.GET("site_msg/conf", middleware.AuthMiddleware, app.UserSiteMessageConfView)
	rg.PUT("site_msg/conf", middleware.AuthMiddleware, app.UserSiteMessageConfUpdateView)
	rg.POST("site_msg", middleware.AuthMiddleware, app.SiteMsgReadView)
	rg.DELETE("site_msg", middleware.AuthMiddleware, app.SiteMsgRemoveView)
	rg.GET("site_msg/user", middleware.AuthMiddleware, app.UserMsgView)
}
