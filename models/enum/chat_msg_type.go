package enum

// ChatMsgType 聊天消息类型。
type ChatMsgType int8

const (
	ChatTextMsgType     ChatMsgType = 1
	ChatImageMsgType    ChatMsgType = 2
	ChatMarkdownMsgType ChatMsgType = 3
	ChatReadMsgType     ChatMsgType = 4
)
