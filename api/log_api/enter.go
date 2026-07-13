package log_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/service/log_service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type LogApi struct {
}
type LogListRequest struct {
	common.PageInfo
	Level       enum.LogLevelType `json:"level"`
	UserID      uint              `json:"userId"`
	Ip          string            `json:"ip"`
	ServiceName string            `gorm:"size:32" json:"serviceName"`
	LogType     enum.LogType      `json:"logType"`
}
type LogListResponse struct {
	models.LogModel
	UserNickname string `json:"userNickname"`
	UserAvatar   string `json:"userAvatar"`
}

func (l *LogApi) GetLogList(c *gin.Context) {
	var req LogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}

	list, count, err := common.ListQuery(models.LogModel{
		LogType:     req.LogType,
		Level:       req.Level,
		UserID:      req.UserID,
		Ip:          req.Ip,
		ServiceName: req.ServiceName,
	}, common.Options{
		PageInfo: req.PageInfo,
		Likes:    []string{"title"},
		Preloads: []string{"UserModel"},
	})
	if err != nil {
		res.FailWithData(err.Error(), "获取日志列表失败", c)
		return
	}

	respList := make([]LogListResponse, 0, len(list))
	for _, logModel := range list {
		respList = append(respList, LogListResponse{
			LogModel:     logModel,
			UserNickname: logModel.UserModel.Nickname,
			UserAvatar:   logModel.UserModel.Avatar,
		})
	}

	res.OkWithList(respList, count, c)
}

func (l *LogApi) LogReadView(c *gin.Context) {
	var req models.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		res.FailWithError(err, c)
		return
	}

	var logModel models.LogModel
	result := global.DB.Take(&logModel, req.ID)
	if result.Error != nil {
		res.FailWithMsg("日志不存在", c)
		return
	}
	if logModel.IsRead {
		res.FailWithMsg("日志已读取", c)
		return
	}

	result = global.DB.Model(&logModel).Update("is_read", true)
	if result.Error != nil {
		res.FailWithMsg("更新日志详情失败", c)
		return
	}
	if result.RowsAffected == 0 {
		res.FailWithMsg("日志读取失败", c)
		return
	}

	res.OkWithMsg("日志读取成功", c)
}
func (l *LogApi) LogRemoveView(c *gin.Context) {
	var req models.RemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	log := log_service.GetLog(c)
	log.ShowRequest()
	log.ShowResponse()

	var logModel []models.LogModel
	result := global.DB.Find(&logModel, req.IDList)
	if result.Error != nil {
		res.FailWithMsg("查询日志失败", c)
		return
	}
	if len(logModel) == 0 {
		res.FailWithMsg("日志不存在", c)
		return
	}

	result = global.DB.Delete(&logModel)
	if result.Error != nil {
		res.FailWithMsg("删除日志失败", c)
		return
	}

	msg := fmt.Sprintf("删除%d条日志成功", result.RowsAffected)
	res.OkWithMsg(msg, c)
}
