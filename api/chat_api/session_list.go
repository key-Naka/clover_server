package chat_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/ctype/chat_msg"
	"clover_server/models/enum/relationship_enum"
	"clover_server/utils/jwts"
	"time"

	"github.com/gin-gonic/gin"
)

type SessionListRequest struct { common.PageInfo }

type SessionListResponse struct {
	UserID       uint                       `json:"userID"`
	UserNickname string                     `json:"userNickname"`
	UserAvatar   string                     `json:"userAvatar"`
	Msg          chat_msg.ChatMsg           `json:"msg"`
	MsgType      int8                       `json:"msgType"`
	NewMsgDate   time.Time                  `json:"newMsgDate"`
	Relation     relationship_enum.Relation `json:"relation"`
}

func (ChatApi) SessionListView(c *gin.Context) {
	var req SessionListRequest
	if err := c.ShouldBindQuery(&req); err != nil { res.FailWithError(err, c); return }
	claims := jwts.GetClaims(c)
	if claims == nil { res.FailWithMsg("请登录", c); return }
	var deletedIDs []uint
	global.DB.Model(&models.UserChatActionModel{}).Where("user_id = ? and is_delete = ?", claims.UserID, true).Pluck("chat_id", &deletedIDs)
	var chats []models.ChatModel
	query := global.DB.Preload("SendUserModel").Preload("RevUserModel").Where("send_user_id = ? or rev_user_id = ?", claims.UserID, claims.UserID).Order("id desc")
	if len(deletedIDs) > 0 { query = query.Where("id not in ?", deletedIDs) }
	if err := query.Find(&chats).Error; err != nil { res.FailWithError(err, c); return }
	seen := map[uint]bool{}
	list := make([]SessionListResponse, 0)
	for _, chat := range chats {
		otherID := chat.SendUserID
		otherUser := chat.SendUserModel
		if chat.SendUserID == claims.UserID {
			otherID = chat.RevUserID
			otherUser = chat.RevUserModel
		}
		if seen[otherID] { continue }
		seen[otherID] = true
		list = append(list, SessionListResponse{UserID: otherID, UserNickname: otherUser.Nickname, UserAvatar: otherUser.Avatar, Msg: chat.Msg, MsgType: int8(chat.MsgType), NewMsgDate: chat.CreatedAt, Relation: relationship_enum.RelationStranger})
	}
	count := len(list)
	offset := req.GetOffset()
	if offset >= count { res.OkWithList([]SessionListResponse{}, count, c); return }
	end := offset + req.GetLimit(); if end > count { end = count }
	res.OkWithList(list[offset:end], count, c)
}
