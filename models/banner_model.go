package models

type BannerModel struct {
	Model
	Title string `gorm:"size:256" json:"title"`
	Cover string `gorm:"size:256" json:"cover"`
	Href  string `gorm:"size:256" json:"href"`
}
