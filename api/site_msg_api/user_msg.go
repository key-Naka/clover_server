package site_msg_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

type UserMsgResponse struct {
	CommentMsgCount int `json:"commentMsgCount"`
	DiggMsgCount    int `json:"diggMsgCount"`
	PrivateMsgCount int `json:"privateMsgCount"`
	SystemMsgCount  int `json:"systemMsgCount"`
}

func (SiteMsgApi) UserMsgView(c *gin.Context) {
	claims := jwts.GetClaims(c)
	if claims == nil { res.FailWithMsg("请登录", c); return }
	var msgList []models.MessageModel
	global.DB.Find(&msgList, "rev_user_id = ? and is_read = ?", claims.UserID, false)
	var data UserMsgResponse
	for _, model := range msgList {
		switch model.Type {
		case enum.MessageCommentType, enum.MessageCommentReplyType:
			data.CommentMsgCount++
		case enum.MessageArticleLikeType, enum.MessageCommentLikeType, enum.MessageArticleCollectType:
			data.DiggMsgCount++
		case enum.MessageSystemNotificationType:
			data.SystemMsgCount++
		}
	}
	var chatList []models.ChatModel
	global.DB.Find(&chatList, "rev_user_id = ?", claims.UserID)
	var chatIDList []uint
	for _, model := range chatList { chatIDList = append(chatIDList, model.ID) }
	var actions []models.UserChatActionModel
	if len(chatIDList) > 0 { global.DB.Where("user_id = ? and chat_id in ?", claims.UserID, chatIDList).Find(&actions) }
	actionMap := map[uint]models.UserChatActionModel{}
	for _, item := range actions { actionMap[item.ChatID] = item }
	for _, model := range chatList { if action, ok := actionMap[model.ID]; !ok || !action.IsRead { data.PrivateMsgCount++ } }
	var records []models.UserGlobalNotificationModel
	global.DB.Where("user_id = ? and (is_read = ? or is_delete = ?)", claims.UserID, true, true).Find(&records)
	excludeIDs := make([]uint, 0, len(records))
	for _, item := range records { excludeIDs = append(excludeIDs, item.NotificationID) }
	query := global.DB.Model(&models.GlobalNotificationModel{})
	if len(excludeIDs) > 0 { query = query.Where("id not in ?", excludeIDs) }
	var total int64
	query.Count(&total)
	data.SystemMsgCount += int(total)
	res.OkWithData(data, c)
}
