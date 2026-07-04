package site_api

import (
	"clover_server/models/enum"
	"clover_server/service/log_service"

	"github.com/gin-gonic/gin"
)

type SiteApi struct {
}

func (s *SiteApi) SiteinfoView(c *gin.Context) {
	log_service.NewLoginSuccess(c, enum.UserPwdLoginType)
	log_service.NewLoginFail(c, enum.UserPwdLoginType, "登录失败", "admin", "123456")
	c.JSON(200, gin.H{
		"message": "站点信息",
	})
}

func (s *SiteApi) SiteUpdateView(c *gin.Context) {
	ActionLog := log_service.GetLog(c)
	ActionLog.Save()
}
