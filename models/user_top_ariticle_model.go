package models

import "time"

type UserTopArticleModel struct {
	UserID       uint         `gorm:"uniqueIndex:idx_name" json:"userId"`
	UserModel    UserModel    `gorm:"foreignKey:UserID " json:"-"`
	ArticleID    uint         `gorm:"uniqueIndex:idx_name" json:"articleId"`
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID" json:"-"`
	CreatedAt    time.Time    `json:"createdAt"`
}
