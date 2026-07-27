package router

import (
	"clover_server/api"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
)

func BannerRouter(r *gin.RouterGroup) {
	app := api.App.BannerApi
	r.GET("banner", app.BannerList)
	r.POST("banner", middleware.AdminAuthMiddleware, app.BannerCreate)
	r.PUT("banner/:id", middleware.AdminAuthMiddleware, app.BannerUpdate)
	r.DELETE("banner", middleware.AdminAuthMiddleware, app.BannerRemove)
}
