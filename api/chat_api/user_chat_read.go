package chat_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/ctype/chat_msg"
	"clover_server/models/enum"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

func (ChatApi) ChatReadView(c *gin.Context) {
	var req models.IDRequest
	if err := c.ShouldBindUri(&req); err != nil { res.FailWithError(err, c); return }
	claims := jwts.GetClaims(c)
	if claims == nil { res.FailWithMsg("请登录", c); return }
	var chat models.ChatModel
	if err := global.DB.Take(&chat, req.ID).Error; err != nil { res.FailWithMsg("聊天记录不存在", c); return }
	var action models.UserChatActionModel
	if err := global.DB.Take(&action, "user_id = ? and chat_id = ?", claims.UserID, chat.ID).Error; err != nil {
		action = models.UserChatActionModel{UserID: claims.UserID, ChatID: chat.ID, IsRead: true}
		if err = global.DB.Create(&action).Error; err != nil { res.FailWithError(err, c); return }
	} else {
		if action.IsDelete { res.FailWithMsg("该记录已删除", c); return }
		if err = global.DB.Model(&action).Update("is_read", true).Error; err != nil { res.FailWithError(err, c); return }
	}
	sendToUser(chat.SendUserID, ChatResponse{ChatListResponse: ChatListResponse{ChatModel: models.ChatModel{Model: models.Model{ID: chat.ID}, SendUserID: claims.UserID, RevUserID: chat.SendUserID, MsgType: enum.ChatReadMsgType, Msg: chat_msg.ChatMsg{MsgReadMsg: &chat_msg.MsgReadMsg{ReadChatID: chat.ID}}}}})
	res.OkWithMsg("已读成功", c)
}
