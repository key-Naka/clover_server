package router

import (
	"clover_server/api"
	"clover_server/middleware"

	"github.com/gin-gonic/gin"
)

func ArticleRouter(r *gin.RouterGroup) {
	app := api.App.ArticleApi

	r.GET("article/list", app.ArticleListView)
	r.GET("article/detail", app.ArticleDetailView)

	r.POST("article", middleware.AuthMiddleware, app.ArticleCreateView)
	r.PUT("article", middleware.AuthMiddleware, app.ArticleUpdateView)
	r.DELETE("article", middleware.AuthMiddleware, app.ArticleRemoveUserView)

	r.GET("article/category/list", app.CategoryListView)
	r.GET("article/category/options", middleware.AuthMiddleware, app.CategoryOptionsView)
	r.POST("article/category", middleware.AuthMiddleware, app.CategoryCreateView)
}
