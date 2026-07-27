// router/site_router.go
package router

import (
	"clover_server/api"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
)

func ImageRouter(r *gin.RouterGroup) {
	app := api.App.ImageApi
	r.POST("images", middleware.AdminAuthMiddleware, app.ImageUploadView)
	r.POST("images/qiniu", middleware.AuthMiddleware, app.QiNiuGenToken)
	// r.POST("images/transfer_deposit", middleware.AuthMiddleware, middleware.BindJsonMiddleware[image_api.TransferDepositRequest], app.TransferDepositView)
	r.GET("images", middleware.AdminAuthMiddleware, app.ImageListView)
	r.DELETE("images", middleware.AdminAuthMiddleware, app.ImageRemoveView)
}
