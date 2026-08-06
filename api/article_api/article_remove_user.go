package article_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/service/es_service"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

// ArticleRemoveUserView 删除当前用户自己的文章。
// @Summary 删除文章
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.IDRequest true "文章 ID"
// @Success 200 {object} res.MessageResponse "文章删除成功"
// @Failure 400 {object} res.ErrorResponse "参数、归属或认证失败"
// @Failure 500 {object} res.ErrorResponse "服务异常"
// @Router /article [delete]
func (ArticleApi) ArticleRemoveUserView(c *gin.Context) {
	var req models.IDRequest
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
		res.FailWithMsg("无权限删除此文章", c)
		return
	}

	if err := global.DB.Delete(&article).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	if err := es_service.DeleteArticleDocument(article.ID); err != nil && global.ES != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}
	res.OkWithMsg("文章删除成功", c)
}
