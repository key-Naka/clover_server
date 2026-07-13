package log_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"

	"github.com/gin-gonic/gin"
)

type LogApi struct {
}
type LogListRequest struct {
	Page        int               `form:"page"`
	Limit       int               `form:"limit"`
	Key         string            `form:"key"`
	Level       enum.LogLevelType `json:"level"`
	UserID      uint              `json:"userId"`
	Ip          string            `json:"ip"`
	ServiceName string            `gorm:"size:32" json:"serviceName"`
	LogType     enum.LogType      `json:"logType"`
}

func (l *LogApi) GetLogList(c *gin.Context) {
	var req LogListRequest
	c.ShouldBindQuery(&req)
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 5
	}
	var List []models.LogModel
	model := &models.LogModel{
		LogType:     req.LogType,
		Level:       req.Level,
		UserID:      req.UserID,
		Ip:          req.Ip,
		ServiceName: req.ServiceName,
	}
	global.DB.Debug().Where(model).Offset((req.Page - 1) * req.Limit).Limit(req.Limit).Find(&List)
	var count int64
	global.DB.Debug().Where(model).Count(&count)
	res.OkWithList(List, int(count), c)
	return
}
