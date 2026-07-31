package comment_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/utils/jwts"
	"strings"

	"github.com/gin-gonic/gin"
)

type CommentCreateRequest struct {
	ArticleID uint  `json:"articleID" binding:"required"`
	ParentID  *uint `json:"parentID"`
	Content   string `json:"content" binding:"required"`
}

func (CommentApi) CommentCreateView(c *gin.Context) {
	var req CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	claims := jwts.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}
	article, err := loadArticle(req.ArticleID)
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}
	if !article.OpenComment {
		res.FailWithMsg("该文章未开启评论", c)
		return
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		res.FailWithMsg("评论内容不能为空", c)
		return
	}
	comment := models.CommentModel{
		Content:   content,
		ArticleID: req.ArticleID,
		UserID:    claims.UserID,
		ParentID:  req.ParentID,
	}
	if req.ParentID != nil {
		var parent models.CommentModel
		if err = global.DB.Take(&parent, "id = ? and article_id = ?", *req.ParentID, req.ArticleID).Error; err != nil {
			res.FailWithMsg("父评论不存在", c)
			return
		}
		comment.RootParentID = findRootParentID(parent)
		if parent.ParentID == nil {
			comment.RootParentID = &parent.ID
		}
	}
	if err = global.DB.Create(&comment).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	global.DB.Model(article).Update("comment_count", article.CommentCount+1)
	if err = syncArticleCommentToES(req.ArticleID); err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}
	res.OkWithData(comment, c)
}
