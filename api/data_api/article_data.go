package data_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"time"

	"github.com/gin-gonic/gin"
)

type MonthData struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type ArticleYearDataResponse struct {
	List       []MonthData `json:"list"`
	Rate       float64     `json:"rate"`
	Increase   int64       `json:"increase"`
	TotalCount int64       `json:"totalCount"`
}

func (DataApi) ArticleYearDataView(c *gin.Context) {
	now := time.Now()
	list := make([]MonthData, 0, 12)
	var total int64
	var prev int64
	var current int64
	for i := 11; i >= 0; i-- {
		month := now.AddDate(0, -i, 0)
		start := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
		end := start.AddDate(0, 1, 0)
		var count int64
		global.DB.Model(&models.ArticleModel{}).Where("status = ? and created_at >= ? and created_at < ?", 3, start, end).Count(&count)
		list = append(list, MonthData{Date: start.Format("2006-01"), Count: count})
		total += count
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
	res.OkWithData(ArticleYearDataResponse{List: list, Rate: rate, Increase: current - prev, TotalCount: total}, c)
}
