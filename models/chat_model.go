package models

import (
	"clover_server/models/ctype/chat_msg"
	"clover_server/models/enum"
)

type ChatModel struct {
	Model
	SendUserID    uint             `json:"sendUserID"`
	SendUserModel UserModel        `gorm:"foreignKey:SendUserID"  json:"-"`
	RevUserID     uint             `json:"revUserID"`
	RevUserModel  UserModel        `gorm:"foreignKey:RevUserID"  json:"-"`
	MsgType       enum.ChatMsgType `json:"msgType"` // 消息类型
	Msg           chat_msg.ChatMsg `gorm:"type:longtext;serializer:json" json:"msg"`
}
