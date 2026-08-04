package chat_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

func (ChatApi) UserChatDeleteView(c *gin.Context) {
	var req models.RemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil { res.FailWithError(err, c); return }
	claims := jwts.GetClaims(c)
	if claims == nil { res.FailWithMsg("请登录", c); return }
	var list []models.UserChatActionModel
	global.DB.Where("user_id = ? and chat_id in ?", claims.UserID, req.IDList).Find(&list)
	existMap := map[uint]models.UserChatActionModel{}
	for _, item := range list { existMap[item.ChatID] = item }
	for _, chatID := range req.IDList {
		item, ok := existMap[chatID]
		if !ok { _ = global.DB.Create(&models.UserChatActionModel{UserID: claims.UserID, ChatID: chatID, IsDelete: true}).Error; continue }
		if !item.IsDelete { _ = global.DB.Model(&item).Update("is_delete", true).Error }
	}
	res.OkWithMsg("删除成功", c)
}
