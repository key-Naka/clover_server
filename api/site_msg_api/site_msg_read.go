package site_msg_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

type SiteMsgReadRequest struct {
	ID uint `json:"id"`
	T  int8 `json:"t" binding:"required,oneof=1 2 3"`
}

func (SiteMsgApi) SiteMsgReadView(c *gin.Context) {
	var req SiteMsgReadRequest
	if err := c.ShouldBindJSON(&req); err != nil { res.FailWithError(err, c); return }
	claims := jwts.GetClaims(c)
	if claims == nil { res.FailWithMsg("请登录", c); return }
	if req.ID != 0 {
		var msg models.MessageModel
		if err := global.DB.Take(&msg, "id = ? and rev_user_id = ?", req.ID, claims.UserID).Error; err != nil { res.FailWithMsg("消息不存在", c); return }
		if msg.IsRead { res.FailWithMsg("消息已读", c); return }
		if err := global.DB.Model(&msg).Update("is_read", true).Error; err != nil { res.FailWithError(err, c); return }
		res.OkWithMsg("消息已读", c); return
	}
	typeList := getSiteMessageTypes(req.T)
	result := global.DB.Model(&models.MessageModel{}).Where("rev_user_id = ? and is_read = ? and type in ?", claims.UserID, false, typeList).Update("is_read", true)
	if result.Error != nil { res.FailWithError(result.Error, c); return }
	res.OkWithMsg("已读成功", c)
}
