package chat_service

import (
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/ctype/chat_msg"
	"clover_server/models/enum"
	"strings"
)

func ToChat(sendUserID, revUserID uint, msgType enum.ChatMsgType, msg chat_msg.ChatMsg) error {
	return global.DB.Create(&models.ChatModel{
		SendUserID: sendUserID,
		RevUserID:  revUserID,
		MsgType:    msgType,
		Msg:        msg,
	}).Error
}

func ToTextChat(sendUserID, revUserID uint, content string) error {
	return ToChat(sendUserID, revUserID, enum.ChatTextMsgType, chat_msg.ChatMsg{TextMsg: &chat_msg.TextMsg{Content: content}})
}

func ToImageChat(sendUserID, revUserID uint, src string) error {
	return ToChat(sendUserID, revUserID, enum.ChatImageMsgType, chat_msg.ChatMsg{ImageMsg: &chat_msg.ImageMsg{Src: src}})
}

func ToMarkdownChat(sendUserID, revUserID uint, content string) error {
	return ToChat(sendUserID, revUserID, enum.ChatMarkdownMsgType, chat_msg.ChatMsg{MarkdownMsg: &chat_msg.MarkdownMsg{Content: strings.TrimSpace(content)}})
}
