package router

import (
	"clover_server/global"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(global.Config.System.GinMode)
	r := gin.Default()
	rg := r.Group("/api/")
	rg.Use(middleware.LogMiddleware)
	ArticleRouter(rg)
	CommentRouter(rg)
	DataRouter(rg)
	FocusRouter(rg)
	GlobalNotificationRouter(rg)
	ImageRouter(rg)
	SiteRouter(rg)
	LogRouter(rg)
	BannerRouter(rg)
	CaptchaRouter(rg)
	UserRouter(rg)

	r.Run(global.Config.System.GetAddr())
}
