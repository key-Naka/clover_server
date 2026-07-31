package data_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"time"

	"github.com/gin-gonic/gin"
)

type GrowthRequest struct {
	Type int8 `form:"type" binding:"required,oneof=1 2 3"`
}

type DayData struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type GrowthDataResponse struct {
	List     []DayData `json:"list"`
	Rate     float64   `json:"rate"`
	Increase int64     `json:"increase"`
}

func (DataApi) GrowthDataView(c *gin.Context) {
	var req GrowthRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	list := make([]DayData, 0, 7)
	now := time.Now()
	var prev int64
	var current int64
	for i := 6; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
		end := start.AddDate(0, 0, 1)
		var count int64
		switch req.Type {
		case 1:
			global.DB.Model(&models.SiteFlowModel{}).Where("created_at >= ? and created_at < ?", start, end).Count(&count)
		case 2:
			global.DB.Model(&models.ArticleModel{}).Where("status = ? and created_at >= ? and created_at < ?", 3, start, end).Count(&count)
		case 3:
			global.DB.Model(&models.UserModel{}).Where("created_at >= ? and created_at < ?", start, end).Count(&count)
		}
		list = append(list, DayData{Date: start.Format("2006-01-02"), Count: count})
		if i == 1 {
			prev = count
		}
		if i == 0 {
			current = count
		}
	}
	rate := 0.0
	if prev > 0 {
		rate = float64(current-prev) / float64(prev) * 100
	}
	res.OkWithData(GrowthDataResponse{List: list, Rate: rate, Increase: current - prev}, c)
}
