package router

import (
	"clover_server/api"

	"github.com/gin-gonic/gin"
)

func SearchRouter(rg *gin.RouterGroup) {
	app := api.App.SearchApi
	rg.GET("search/article", app.ArticleSearchView)
	rg.GET("search/tags", app.TagAggView)
	rg.GET("search/text", app.TextSearchView)
}
