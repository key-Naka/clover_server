package chat_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/ctype/chat_msg"
	"clover_server/models/enum"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

type ChatApi struct{}

type ChatListRequest struct {
	common.PageInfo
	SendUserID uint `form:"sendUserID"`
	RevUserID  uint `form:"revUserID" binding:"required"`
	Type       int8 `form:"type" binding:"required,oneof=1 2"`
}

type ChatListResponse struct {
	models.ChatModel
	SendUserNickname string `json:"sendUserNickname"`
	SendUserAvatar   string `json:"sendUserAvatar"`
	RevUserNickname  string `json:"revUserNickname"`
	RevUserAvatar    string `json:"revUserAvatar"`
	IsMe             bool   `json:"isMe"`
	IsRead           bool   `json:"isRead"`
}

func (ChatApi) ChatListView(c *gin.Context) {
	var req ChatListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	claims := jwts.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}
	query := global.DB.Model(&models.ChatModel{}).Preload("SendUserModel").Preload("RevUserModel")
	if req.Type == 1 {
		query = query.Where("((send_user_id = ? and rev_user_id = ?) or (send_user_id = ? and rev_user_id = ?))", claims.UserID, req.RevUserID, req.RevUserID, claims.UserID)
	} else {
		if claims.Role != enum.AdminRole {
			res.FailWithMsg("角色错误", c)
			return
		}
		if req.SendUserID == 0 {
			res.FailWithMsg("发送人不能为空", c)
			return
		}
		query = query.Where("((send_user_id = ? and rev_user_id = ?) or (send_user_id = ? and rev_user_id = ?))", req.SendUserID, req.RevUserID, req.RevUserID, req.SendUserID)
	}
	var deletedIDs []uint
	global.DB.Model(&models.UserChatActionModel{}).Where("user_id = ? and is_delete = ?", claims.UserID, true).Pluck("chat_id", &deletedIDs)
	if len(deletedIDs) > 0 {
		query = query.Where("id not in ?", deletedIDs)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	var list []models.ChatModel
	if err := query.Order("id desc").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&list).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	var chatIDs []uint
	for _, item := range list {
		chatIDs = append(chatIDs, item.ID)
	}
	var actions []models.UserChatActionModel
	if len(chatIDs) > 0 {
		global.DB.Where("user_id = ? and chat_id in ?", claims.UserID, chatIDs).Find(&actions)
	}
	readMap := map[uint]bool{}
	for _, item := range actions {
		readMap[item.ChatID] = item.IsRead
	}
	responseList := make([]ChatListResponse, 0, len(list))
	for _, item := range list {
		responseList = append(responseList, ChatListResponse{
			ChatModel:        item,
			SendUserNickname: item.SendUserModel.Nickname,
			SendUserAvatar:   item.SendUserModel.Avatar,
			RevUserNickname:  item.RevUserModel.Nickname,
			RevUserAvatar:    item.RevUserModel.Avatar,
			IsMe:             item.SendUserID == claims.UserID,
			IsRead:           readMap[item.ID],
		})
	}
	res.OkWithList(responseList, int(count), c)
}

type ChatRequest struct {
	RevUserID uint             `json:"revUserID"`
	MsgType   enum.ChatMsgType `json:"msgType"`
	Msg       chat_msg.ChatMsg `json:"msg"`
}
