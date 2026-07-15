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

// LogReadView 是一个处理日志读取视图的函数，它接收一个 gin.Context 类型的参数
func (l *LogApi) LogReadView(c *gin.Context) {
    // 定义一个 models.IDRequest 类型的变量 req，用于接收请求参数
	var req models.IDRequest
    // 使用 ShouldBindUri 方法将 URI 参数绑定到 req 结构体上，如果绑定失败则返回错误
	if err := c.ShouldBindUri(&req); err != nil {
		res.FailWithError(err, c)
		return
	}

    // 定义一个 models.LogModel 类型的变量 logModel，用于存储日志数据
	var logModel models.LogModel
    // 使用 global.DB.Take 方法根据 ID 查询日志数据，如果查询失败则返回错误
	result := global.DB.Take(&logModel, req.ID)
	if result.Error != nil {
		res.FailWithMsg("日志不存在", c)
		return
	}
    // 检查日志是否已被读取，如果已读取则返回错误
	if logModel.IsRead {
		res.FailWithMsg("日志已读取", c)
		return
	}

    // 使用 global.DB.Model 更新日志的 is_read 字段为 true，表示已读取
	result = global.DB.Model(&logModel).Update("is_read", true)
	if result.Error != nil {
		res.FailWithMsg("更新日志详情失败", c)
		return
	}
    // 检查是否有行被更新，如果没有则说明更新失败
	if result.RowsAffected == 0 {
		res.FailWithMsg("日志读取失败", c)
		return
	}

    // 如果所有操作都成功，则返回成功消息
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
