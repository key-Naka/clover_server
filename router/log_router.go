package router

import (
	"clover_server/api"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
)

func LogRouter(r *gin.RouterGroup) {
	App := api.App.LogApi

	r.GET("logs", middleware.AuthMiddleware, App.GetLogList)
	r.GET("logs/:id", middleware.AuthMiddleware, App.LogReadView)
	r.DELETE("logs", middleware.AdminAuthMiddleware, App.LogRemoveView)
}
