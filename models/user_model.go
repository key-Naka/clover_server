package models

type UserModel struct {
	Model
	Username       string   `gorm:"size:32" json:Username`
	Nickname       string   `gorm:"size:32"  json:Nickname`
	Avatar         string   `gorm:"size:256"  json:Avatar`
	Abstract       string   `gorm:"size:256"  json:Abstract`
	RegisterSource int8     ` json:"RegisterSource"`
	Password       string   `gorm:"size:256"  json:password`
	OpenID string `gorm:"size:64"  json:openID`
}
