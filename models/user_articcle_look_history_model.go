package models

type UserArticleLookHistoryModel struct {
	Model
	UserID       uint         `json:"userId"`
	UserModel    UserModel    `gorm:"foreignKey:UserID " json:"-"`
	ArticleID    uint         `json:"articleId"`
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID" json:"-"`
}
