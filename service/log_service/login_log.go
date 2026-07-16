package log_service

import (
	"clover_server/core"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func NewLoginSuccess(c *gin.Context, loginType enum.LoginType) {
	ip := c.ClientIP()
	region := core.SearchAddr(ip)
	// token := c.GetHeader("token")
	userID := uint(1)
	username := "admin"
	logModel := models.LogModel{
		LogType:     enum.LoginLogType,
		LoginType:   loginType,
		Title:       "用户登录",
		Content:     "登录成功",
		UserID:      userID,
		Ip:          ip,
		Addr:        region,
		LoginStatus: true,
		Username:    username,
		Password:    "123456",
	}
	if err := global.DB.Create(&logModel).Error; err != nil {
		slog.Error("写入登录成功日志失败", "err", err, "login_type", loginType, "user_id", userID, "username", username, "client_ip", ip)
		return
	}

	slog.Info("用户登录成功", "login_type", loginType, "user_id", userID, "username", username, "client_ip", ip, "region", region)
}
func NewLoginFail(c *gin.Context, loginType enum.LoginType, msg string, username string, password string) {
	ip := c.ClientIP()
	region := core.SearchAddr(ip)

	logModel := models.LogModel{
		LogType:     enum.LoginLogType,
		LoginType:   loginType,
		Title:       "用户登录失败",
		Content:     msg,
		Ip:          ip,
		Addr:        region,
		LoginStatus: false,
		Username:    username,
		Password:    password,
	}
	if err := global.DB.Create(&logModel).Error; err != nil {
		slog.Error("写入登录失败日志失败", "err", err, "login_type", loginType, "username", username, "client_ip", ip, "reason", msg)
		return
	}

	slog.Warn("用户登录失败", "login_type", loginType, "username", username, "client_ip", ip, "region", region, "reason", msg)
}
