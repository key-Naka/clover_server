package models

import (
	"time"
)

type UserModel struct {
	Model
	Username       string `gorm:"size:32" json:username`
	Nickname       string `gorm:"size:32"  json:nickname`
	Avatar         string `gorm:"size:256"  json:avatar`
	Abstract       string `gorm:"size:256"  json:abstract`
	RegisterSource int8   `json:"registerSource"`
	Password       string `gorm:"size:256" json:"password"`
	Email          string `gorm:"size:128" json:"email"`
	OpenID         string `gorm:"size:64" json:"openID"`
	Role           int8   ` json:"role"` //1管理员 2普通用户 3访客
}

type UserConfModel struct {
	UserID             uint       `gorm:"primaryKey;unique" json:"userID"`
	UserModel          UserModel  `gorm:"foreignKey:UserID" json:"-"`
	LikeTags           []string   `gorm:"type:json;serializer:json" json:"likeTags"`
	UpdateUsernameDate *time.Time `json:"updateUsernameDate"` // 上次修改用户名的时间
	OpenCollect        bool       `json:"openCollect"`        // 公开我的收藏
	OpenFollow         bool       `json:"openFollow"`         // 公开我的关注
	OpenFans           bool       `json:"openFans"`           // 公开我的粉丝
	HomeStyleID        uint       `json:"homeStyleID"`        // 主页样式的id
	LookCount          int        `json:"lookCount"`          // 主页的访问次数
}
