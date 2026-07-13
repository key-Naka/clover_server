package router

import (
	"clover_server/api"

	"github.com/gin-gonic/gin"
)

func SiteRouter(r *gin.RouterGroup) {
	App := api.App.SiteApi

	r.GET("site", App.SiteinfoView)
	r.POST("site", App.SiteUpdateView)
}
