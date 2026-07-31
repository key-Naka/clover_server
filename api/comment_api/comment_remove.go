package comment_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

func (CommentApi) CommentRemoveView(c *gin.Context) {
	var req models.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	claims := jwts.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}
	var comment models.CommentModel
	if err := global.DB.Take(&comment, req.ID).Error; err != nil {
		res.FailWithMsg("评论不存在", c)
		return
	}
	article, err := loadArticle(comment.ArticleID)
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}
	if claims.Role != enum.AdminRole && claims.UserID != comment.UserID && claims.UserID != article.UserID {
		res.FailWithMsg("无权限删除该评论", c)
		return
	}
	idList, err := collectCommentIDs(comment.ID)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	if err = global.DB.Where("id in ?", idList).Delete(&models.CommentModel{}).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	if article.CommentCount >= len(idList) {
		global.DB.Model(article).Update("comment_count", article.CommentCount-len(idList))
	}
	if err = syncArticleCommentToES(comment.ArticleID); err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}
	res.OkWithMsg("评论删除成功", c)
}
