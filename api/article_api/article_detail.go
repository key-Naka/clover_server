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
