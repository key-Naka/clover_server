package router

import (
	"clover_server/api"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
)

func DataRouter(rg *gin.RouterGroup) {
	app := api.App.DataApi
	rg.GET("data/sum", middleware.AdminAuthMiddleware, app.SumView)
	rg.GET("data/article/year", middleware.AdminAuthMiddleware, app.ArticleYearDataView)
	rg.GET("data/computer", middleware.AdminAuthMiddleware, app.ComputerDataView)
	rg.GET("data/growth", middleware.AdminAuthMiddleware, app.GrowthDataView)
}
