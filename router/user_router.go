package router

import (
	"clover_server/api"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(rg *gin.Engine) {
	var App = api.App.UserApi
	rg.POST("user/send_email", middleware.EmailVerifyMiddleware, App.SendEmail)
	rg.POST("user/email", middleware.EmailVerifyMiddleware, App.RegisterEmailView)
	rg.POST("user/qq", App.QQLoginView)
	rg.POST("user/login", middleware.CaptchaMiddleware, App.PwdLoginApi)
}
