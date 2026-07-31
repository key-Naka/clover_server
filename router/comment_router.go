package router

import (
	"clover_server/api"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
)

func CommentRouter(rg *gin.RouterGroup) {
	app := api.App.CommentApi
	rg.POST("comment", middleware.AuthMiddleware, app.CommentCreateView)
	rg.GET("comment/tree/:id", app.CommentTreeView)
	rg.GET("comment", middleware.AuthMiddleware, app.CommentListView)
	rg.DELETE("comment/:id", middleware.AuthMiddleware, app.CommentRemoveView)
	rg.GET("comment/digg/:id", middleware.AuthMiddleware, app.CommentDiggView)
}
