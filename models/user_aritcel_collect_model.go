package models

import "time"

type UserArticleCollectModel struct {
	CreatedAt time.Time `json:"createdAt"`

	UserID       uint         `gorm:"uniqueIndex:idx_name" json:"userId"`
	UserModel    UserModel    `gorm:"foreignKey:UserID " json:"-"`
	ArticleID    uint         `gorm:"uniqueIndex:idx_name" json:"articleId"`
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID" json:"-"`
	CollectID    uint         `gorm:"uniqueIndex:idx_name" json:"collectId"`
	CollectModel CollectModel `gorm:"foreignKey:CollectID" json:"-"`
}
