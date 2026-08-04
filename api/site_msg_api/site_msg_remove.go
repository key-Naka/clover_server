package site_msg_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

type SiteMsgRemoveRequest struct {
	ID uint `json:"id"`
	T  int8 `json:"t" binding:"required,oneof=1 2 3"`
}

func (SiteMsgApi) SiteMsgRemoveView(c *gin.Context) {
	var req SiteMsgRemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil { res.FailWithError(err, c); return }
	claims := jwts.GetClaims(c)
	if claims == nil { res.FailWithMsg("请登录", c); return }
	if req.ID != 0 {
		if err := global.DB.Delete(&models.MessageModel{}, "id = ? and rev_user_id = ?", req.ID, claims.UserID).Error; err != nil { res.FailWithError(err, c); return }
		res.OkWithMsg("删除成功", c); return
	}
	typeList := getSiteMessageTypes(req.T)
	if err := global.DB.Delete(&models.MessageModel{}, "rev_user_id = ? and type in ?", claims.UserID, typeList).Error; err != nil { res.FailWithError(err, c); return }
	res.OkWithMsg("删除成功", c)
}
