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
	rg.GET("user/detail", middleware.AuthMiddleware, App.DetailView)
	rg.GET("user/base_info", middleware.AuthMiddleware, App.BaseInfoView)
	rg.POST("user/login_list", middleware.AuthMiddleware, App.UserLoginListView)
	rg.PUT("user/password", middleware.AuthMiddleware, App.UpdatePasswordView)
	rg.PUT("user/password/reset", middleware.EmailVerifyMiddleware, App.ResetPasswordView)
	rg.PUT("user/email/bind", middleware.EmailVerifyMiddleware, middleware.AuthMiddleware, App.BindEmailView)
}
