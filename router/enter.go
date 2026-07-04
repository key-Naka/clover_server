package router

import (
	"clover_server/global"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	nr := r.Group("/api")
	nr.Use(middleware.LogMiddleware)
	SiteRouter(nr)
	r.Run(global.Config.System.GetAddr())
}
