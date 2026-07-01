package log_service

import (
	"clover_server/core"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"

	"github.com/gin-gonic/gin"
)

func NewLoginSuccess(c *gin.Context, loginType enum.LoginType) {
	ip := c.ClientIP()
	region := core.SearchAddr(ip)
	// token := c.GetHeader("token")
	userID := uint(1)
	username := "admin"
	global.DB.Create(&models.LogModel{
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
	})

}
func NewLoginFail(c *gin.Context, loginType enum.LoginType, msg string, username string, password string) {
	ip := c.ClientIP()
	region := core.SearchAddr(ip)

	global.DB.Create(&models.LogModel{
		LogType:     enum.LoginLogType,
		LoginType:   loginType,
		Title:       "用户登录失败",
		Content:     msg,
		Ip:          ip,
		Addr:        region,
		LoginStatus: false,
		Username:    username,
		Password:    password,
	})
}
