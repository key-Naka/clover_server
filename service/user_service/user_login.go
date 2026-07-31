// service/user_service/user_login.go
package user_service

import (
	"clover_server/core"
	"clover_server/global"
	"clover_server/models"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func (u UserService) UserLogin(c *gin.Context) {
	ip := c.ClientIP()
	addr := core.SearchAddr(ip)
	ua := c.GetHeader("User-Agent")
	err := global.DB.Create(&models.UserLoginModel{
		UserID: u.userModel.ID,
		IP:     ip,
		Addr:   addr,
		UA:     ua,
	}).Error
	if err != nil {
		slog.Error("用户登陆日志写入失败", "err", err)
	}
}
