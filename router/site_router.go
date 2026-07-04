package router

import (
	"clover_server/api"

	"github.com/gin-gonic/gin"
)

func SiteRouter(r *gin.RouterGroup) {
	App := api.App

	r.GET("/site", App.SiteApi.SiteinfoView)
	r.POST("/site", App.SiteApi.SiteUpdateView)
}
