package log_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/service/log_service"
	"fmt"
	"log/slog"

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
		slog.Warn("绑定日志列表查询参数失败", "path", c.Request.URL.Path, "client_ip", c.ClientIP(), "err", err)
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
		slog.Error("获取日志列表失败", "err", err, "log_type", req.LogType, "level", req.Level, "user_id", req.UserID, "ip", req.Ip, "service_name", req.ServiceName)
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
	var req models.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		slog.Warn("绑定日志读取参数失败", "path", c.Request.URL.Path, "client_ip", c.ClientIP(), "err", err)
		res.FailWithError(err, c)
		return
	}

	var logModel models.LogModel
	result := global.DB.Take(&logModel, req.ID)
	if result.Error != nil {
		slog.Warn("日志读取失败，日志不存在", "log_id", req.ID)
		res.FailWithMsg("日志不存在", c)
		return
	}
	if logModel.IsRead {
		slog.Info("日志重复读取", "log_id", req.ID)
		res.FailWithMsg("日志已读取", c)
		return
	}

	result = global.DB.Model(&logModel).Update("is_read", true)
	if result.Error != nil {
		slog.Error("更新日志已读状态失败", "err", result.Error, "log_id", req.ID)
		res.FailWithMsg("更新日志详情失败", c)
		return
	}
	if result.RowsAffected == 0 {
		slog.Warn("日志已读状态未更新", "log_id", req.ID)
		res.FailWithMsg("日志读取失败", c)
		return
	}

	res.OkWithMsg("日志读取成功", c)
}
func (l *LogApi) LogRemoveView(c *gin.Context) {
	var req models.RemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("绑定删除日志参数失败", "path", c.Request.URL.Path, "client_ip", c.ClientIP(), "err", err)
		res.FailWithError(err, c)
		return
	}
	log := log_service.GetLog(c)
	log.ShowRequest()
	log.ShowResponse()
	log.SetTitle("删除日志")
	log.SetItemInfo("日志ID列表", req.IDList)

	var logModel []models.LogModel
	result := global.DB.Find(&logModel, req.IDList)
	if result.Error != nil {
		slog.Error("查询待删除日志失败", "err", result.Error, "id_list", req.IDList)
		res.FailWithMsg("查询日志失败", c)
		return
	}
	if len(logModel) == 0 {
		slog.Warn("删除日志失败，日志不存在", "id_list", req.IDList)
		res.FailWithMsg("日志不存在", c)
		return
	}

	result = global.DB.Delete(&logModel)
	if result.Error != nil {
		slog.Error("删除日志失败", "err", result.Error, "id_list", req.IDList)
		res.FailWithMsg("删除日志失败", c)
		return
	}

	log.SetItemInfo("删除数量", result.RowsAffected)
	slog.Info("删除日志成功", "count", result.RowsAffected, "id_list", req.IDList)
	msg := fmt.Sprintf("删除%d条日志成功", result.RowsAffected)
	res.OkWithMsg(msg, c)
}
