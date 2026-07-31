package enum

// MessageType 站内消息类型。
type MessageType int8

const (
	MessageCommentType            MessageType = 1
	MessageCommentReplyType       MessageType = 2
	MessageArticleLikeType        MessageType = 3
	MessageCommentLikeType        MessageType = 4
	MessageArticleCollectType     MessageType = 5
	MessageUserFocusType          MessageType = 6
	MessageSystemNotificationType MessageType = 7
)
