package models

import "time"

type ArticleDiggModel struct {
	CreatedAt time.Time    `json:"createdAt"`
	ArticleID uint         `gorm:"uniqueIndex:idx_article_user" json:"articleId"`
	Article   ArticleModel `gorm:"foreignKey:ArticleID" json:"-"`
	UserID    uint         `gorm:"uniqueIndex:idx_article_user" json:"userId"`
	User      UserModel    `gorm:"foreignKey:UserID" json:"-"`
}
