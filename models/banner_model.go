package models

type BannerModel struct {
	Model
	Show  bool   `gorm:"default:false" json:"show"`
	Title string `gorm:"size:256" json:"title"`
	Cover string `gorm:"size:256" json:"cover"`
	Href  string `gorm:"size:256" json:"href"`
}
