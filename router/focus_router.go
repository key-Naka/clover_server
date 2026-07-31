package router

import (
	"clover_server/api"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
)

func FocusRouter(rg *gin.RouterGroup) {
	app := api.App.FocusApi
	rg.POST("focus", middleware.AuthMiddleware, app.FocusUserView)
	rg.DELETE("focus", middleware.AuthMiddleware, app.UnFocusUserView)
	rg.GET("focus/my_focus", app.FocusUserListView)
	rg.GET("focus/my_fans", app.FansUserListView)
}
