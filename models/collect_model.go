package models

type CollectModel struct {
	Model
	UserID       uint      `gorm:"index" json:"userId"`
	User         UserModel `gorm:"foreignKey:UserID" json:"-"`
	Title        string    `gorm:"size:32" json:"title"`
	Abstract     string    `gorm:"size:256" json:"abstract"`
	Cover        string    `gorm:"size:256" json:"cover"`
	CollectCount int       `json:"collectCount"`
}
