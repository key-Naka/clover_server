package router

import (
	"clover_server/api"

	"github.com/gin-gonic/gin"
)

func SiteRouter(r *gin.RouterGroup) {
	App := api.App.SiteApi

	r.GET("site/:name", App.SiteinfoView)
	r.GET("site/qq_url", App.SiteInfoQQView)
	r.POST("site/:name", App.SiteUpdateView)
}
