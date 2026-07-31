package article_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/utils/jwts"
	"strings"

	"github.com/gin-gonic/gin"
)

type ArticleUpdateRequest struct {
	ID          uint     `json:"id" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	Abstract    string   `json:"abstract"`
	Content     string   `json:"content" binding:"required"`
	CategoryID  *uint    `json:"categoryID"`
	TagList     []string `json:"tagList"`
	Cover       string   `json:"cover"`
	OpenComment bool     `json:"openComment"`
	Status      int8     `json:"status" binding:"required,oneof=1 2 3 4"`
}

func (ArticleApi) ArticleUpdateView(c *gin.Context) {
	var req ArticleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	claims := jwts.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}

	var article models.ArticleModel
	if err := global.DB.Take(&article, req.ID).Error; err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}
	if article.UserID != claims.UserID {
		res.FailWithMsg("无权限修改此文章", c)
		return
	}

	if req.CategoryID != nil {
		var category models.CategoryModel
		if err := global.DB.Take(&category, "id = ? and user_id = ?", *req.CategoryID, claims.UserID).Error; err != nil {
			res.FailWithMsg("文章分类不存在", c)
			return
		}
	}

	abstract := strings.TrimSpace(req.Abstract)
	if abstract == "" {
		abstract = buildArticleAbstract(req.Content)
	}

	updateData := map[string]any{
		"title":        req.Title,
		"abstract":     abstract,
		"content":      req.Content,
		"category_id":  req.CategoryID,
		"tag_list":     req.TagList,
		"cover":        req.Cover,
		"open_comment": req.OpenComment,
		"status":       req.Status,
	}
	if err := global.DB.Model(&article).Updates(updateData).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	if err := syncArticleAfterWrite(article.ID); err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}
	res.OkWithMsg("文章更新成功", c)
}
