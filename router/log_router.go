package router

import (
	"clover_server/api"

	"github.com/gin-gonic/gin"
)

func LogRouter(r *gin.RouterGroup) {
	App := api.App.LogApi

	r.GET("logs", App.GetLogList)
	r.GET("logs/:id", App.LogReadView)
	r.DELETE("logs", App.LogRemoveView)
}
