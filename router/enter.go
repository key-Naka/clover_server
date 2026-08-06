package router

import (
	swaggerDocs "clover_server/docs"
	"clover_server/global"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run() {
	gin.SetMode(global.Config.System.GinMode)
	r := gin.Default()
	swaggerDocs.SwaggerInfo.BasePath = "/api"
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rg := r.Group("/api/")
	rg.Use(middleware.LogMiddleware)
	ArticleRouter(rg)
	CommentRouter(rg)
	DataRouter(rg)
	FocusRouter(rg)
	GlobalNotificationRouter(rg)
	SearchRouter(rg)
	ChatRouter(rg)
	SiteMsgRouter(rg)
	ImageRouter(rg)
	SiteRouter(rg)
	LogRouter(rg)
	BannerRouter(rg)
	CaptchaRouter(rg)
	UserRouter(rg)

	r.Run(global.Config.System.GetAddr())
}
