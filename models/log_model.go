package models

import (
	"clover_server/models/enum"
)

type LogModel struct {
	Model
	LogType     enum.LogType      `json:"logType"`
	Title       string            `json:"title"`
	Content     string            `json:"content"`
	Level       enum.LogLevelType `json:"level"`
	UserID      uint              `json:"userId"`
	UserModel   UserModel         `gorm:"foreignKey:UserID" json:"-"`
	Ip          string            `json:"ip"`
	Addr        string            `json:"addr"`
	IsRead      bool              `json:"isRead"`
	LoginStatus bool              `json:"loginStatus"`
	Username    string            `gorm:"size:32"json:"username"`
	Password    string            `gorm:"size:32"json:"password"`
	LoginType   enum.LoginType    `json:"loginType"`
	ServiceName string            `gorm:"size:32" json:"serviceName"`
}
