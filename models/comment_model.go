package models

type CommentModel struct {
	Model
	Content         string          `gorm:"size:256;type:text" json:"content"`
	UserID          uint            `json:"userId"`
	UserModel       UserModel       `gorm:"foreignKey:UserID" json:"-"`
	ArticleID       uint            `json:"articleId"`
	ParentID        *uint           `json:"parentId"`
	ParentModel     *CommentModel   `gorm:"foreignKey:ParentID" json:"-"`
	RootParentID    *uint           `json:"rootParentId"`
	SubCommentCount []*CommentModel `gorm:"foreignKey:ParentID" json:"-"`
	DiggCount       int             `json:"diggCount"`
	ArticleModel    ArticleModel    `gorm:"foreignKey:ArticleID" json:"-"`
}
