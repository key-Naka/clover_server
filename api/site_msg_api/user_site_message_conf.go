package site_msg_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/utils/jwts"
	"clover_server/utils/mps"

	"github.com/gin-gonic/gin"
)

type UserMessageConfUpdateRequest struct {
	OpenCommentMessage *bool `json:"openCommentMessage" u:"open_comment_message"`
	OpenDiggMessage    *bool `json:"openDiggMessage" u:"open_digg_message"`
	OpenPrivateChat    *bool `json:"openPrivateChat" u:"open_private_chat"`
}

func (SiteMsgApi) UserSiteMessageConfView(c *gin.Context) {
	claims := jwts.GetClaims(c)
	if claims == nil { res.FailWithMsg("请登录", c); return }
	var conf models.UserMessageConfModel
	if err := global.DB.Take(&conf, "user_id = ?", claims.UserID).Error; err != nil { res.FailWithMsg("用户消息配置不存在", c); return }
	res.OkWithData(conf, c)
}

func (SiteMsgApi) UserSiteMessageConfUpdateView(c *gin.Context) {
	var req UserMessageConfUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil { res.FailWithError(err, c); return }
	claims := jwts.GetClaims(c)
	if claims == nil { res.FailWithMsg("请登录", c); return }
	var conf models.UserMessageConfModel
	if err := global.DB.Take(&conf, "user_id = ?", claims.UserID).Error; err != nil { res.FailWithMsg("用户消息配置不存在", c); return }
	updateMap := mps.StructToMap(req, "u")
	if len(updateMap) == 0 { res.OkWithMsg("无可更新字段", c); return }
	if err := global.DB.Model(&conf).Updates(updateMap).Error; err != nil { res.FailWithError(err, c); return }
	res.OkWithMsg("更新成功", c)
}
