package chat_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

func (ChatApi) UserChatDeleteByUserView(c *gin.Context) {
	var req models.IDRequest
	if err := c.ShouldBindUri(&req); err != nil { res.FailWithError(err, c); return }
	claims := jwts.GetClaims(c)
	if claims == nil { res.FailWithMsg("请登录", c); return }
	var user models.UserModel
	if err := global.DB.Take(&user, req.ID).Error; err != nil { res.FailWithMsg("用户不存在", c); return }
	var chats []models.ChatModel
	if err := global.DB.Where("(send_user_id = ? and rev_user_id = ?) or (send_user_id = ? and rev_user_id = ?)", claims.UserID, req.ID, req.ID, claims.UserID).Find(&chats).Error; err != nil { res.FailWithError(err, c); return }
	ids := make([]uint, 0, len(chats))
	for _, item := range chats { ids = append(ids, item.ID) }
	var list []models.UserChatActionModel
	global.DB.Where("user_id = ? and chat_id in ?", claims.UserID, ids).Find(&list)
	existMap := map[uint]models.UserChatActionModel{}
	for _, item := range list { existMap[item.ChatID] = item }
	for _, chatID := range ids {
		item, ok := existMap[chatID]
		if !ok { _ = global.DB.Create(&models.UserChatActionModel{UserID: claims.UserID, ChatID: chatID, IsDelete: true}).Error; continue }
		if !item.IsDelete { _ = global.DB.Model(&item).Update("is_delete", true).Error }
	}
	res.OkWithMsg("删除成功", c)
}
