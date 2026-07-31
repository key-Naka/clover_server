package comment_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

func (CommentApi) CommentDiggView(c *gin.Context) {
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
	var digg models.CommentDiggModel
	err := global.DB.Take(&digg, "user_id = ? and comment_id = ?", claims.UserID, comment.ID).Error
	if err == nil {
		if err = global.DB.Delete(&digg).Error; err != nil {
			res.FailWithError(err, c)
			return
		}
		res.OkWithMsg("取消点赞成功", c)
		return
	}
	digg = models.CommentDiggModel{UserID: claims.UserID, CommentID: comment.ID}
	if err = global.DB.Create(&digg).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg("点赞成功", c)
}
