package models

import "clover_server/models/enum"

type MessageModel struct {
	Model
	Type               enum.MessageType `json:"type"`
	RevUserID          uint             `json:"revUserID"` // 接收人的id
	ActionUserID       uint             `json:"ActionUserID"`
	ActionUserNickname string           `json:"actionUserNickname"`
	ActionUserAvatar   string           `json:"actionUserAvatar"`
	Title              string           `json:"title"`
	Content            string           `json:"content"`
	ArticleID          uint             `json:"articleID"`
	ArticleTitle       string           `json:"articleTitle"`
	CommentID          uint             `json:"commentID"`
	LinkTitle          string           `json:"linkTitle"`
	LinkHref           string           `json:"linkHref"`
	IsRead             bool             `json:"isRead"`
}
