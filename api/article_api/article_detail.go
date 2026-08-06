package article_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

type ArticleDetailResponse struct {
	models.ArticleModel
	Username      string `json:"username"`
	Nickname      string `json:"nickname"`
	UserAvatar    string `json:"userAvatar"`
	CategoryTitle string `json:"categoryTitle,omitempty"`
	IsDigg        bool   `json:"isDigg"`
	IsCollect     bool   `json:"isCollect"`
}

// ArticleDetailView 获取文章详情。
// @Summary 获取文章详情
// @Description 未登录用户只能读取已发布文章；携带 Bearer JWT 时可读取本人文章，管理员可读取任意文章。
// @Tags 文章
// @Produce json
// @Param id query int true "文章 ID"
// @Success 200 {object} res.ArticleDetailResponse "文章详情"
// @Failure 400 {object} res.ErrorResponse "参数或文章可见性校验失败"
// @Failure 500 {object} res.ErrorResponse "服务异常"
// @Router /article/detail [get]
func (ArticleApi) ArticleDetailView(c *gin.Context) {
	var req models.IDRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}

	var article models.ArticleModel
	if err := global.DB.Preload("UserModel").Take(&article, req.ID).Error; err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	claims, _ := jwts.ParseTokenByGin(c)
	if claims == nil {
		if article.Status != 3 {
			res.FailWithMsg("文章不存在", c)
			return
		}
	} else if claims.Role != enum.AdminRole && claims.UserID != article.UserID && article.Status != 3 {
		res.FailWithMsg("文章不存在", c)
		return
	}

	data := ArticleDetailResponse{
		ArticleModel:  article,
		Username:      article.UserModel.Username,
		Nickname:      article.UserModel.Nickname,
		UserAvatar:    article.UserModel.Avatar,
		CategoryTitle: queryCategoryTitle(article.CategoryID),
	}
	if claims != nil {
		var digg models.ArticleDiggModel
		if err := global.DB.Take(&digg, "user_id = ? and article_id = ?", claims.UserID, article.ID).Error; err == nil {
			data.IsDigg = true
		}
		var collect models.UserArticleCollectModel
		if err := global.DB.Take(&collect, "user_id = ? and article_id = ?", claims.UserID, article.ID).Error; err == nil {
			data.IsCollect = true
		}
	}
	res.OkWithData(data, c)
}
