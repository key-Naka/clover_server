package data_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"time"

	"github.com/gin-gonic/gin"
)

type SumResponse struct {
	FlowCount      int64 `json:"flowCount"`
	UserCount      int64 `json:"userCount"`
	ArticleCount   int64 `json:"articleCount"`
	ChatCount      int64 `json:"chatCount"`
	CommentCount   int64 `json:"commentCount"`
	TodayLogin     int64 `json:"todayLogin"`
	TodayRegister  int64 `json:"todayRegister"`
}

func (DataApi) SumView(c *gin.Context) {
	var data SumResponse
	global.DB.Model(&models.SiteFlowModel{}).Count(&data.FlowCount)
	global.DB.Model(&models.UserModel{}).Count(&data.UserCount)
	global.DB.Model(&models.ArticleModel{}).Where("status = ?", 3).Count(&data.ArticleCount)
	global.DB.Model(&models.ChatModel{}).Count(&data.ChatCount)
	global.DB.Model(&models.CommentModel{}).Count(&data.CommentCount)
	start := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	global.DB.Model(&models.UserLoginModel{}).Where("created_at >= ?", start).Count(&data.TodayLogin)
	global.DB.Model(&models.UserModel{}).Where("created_at >= ?", start).Count(&data.TodayRegister)
	res.OkWithData(data, c)
}
