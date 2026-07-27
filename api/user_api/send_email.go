package user_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/service/email_service"
	"clover_server/utils/email_store"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type SendEmailRequest struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}
type SendEmailResponse struct {
	EmailID string `json:"emailID"`
}

func (u *UserApi) SendEmail(c *gin.Context) {
	var cr SendEmailRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		slog.Warn("绑定发送邮箱请求参数失败", "path", c.Request.URL.Path, "client_ip", c.ClientIP(), "err", err)
		res.FailWithMsg(err.Error(), c)
		return
	}
	if !global.Config.Site.Login.EmailLogin {
		res.FailWithMsg("站点未启用邮箱注册", c)
		return
	}
	code := base64Captcha.RandText(4, "0123456789")
	id := base64Captcha.RandomId()
	emailType, err := strconv.Atoi(cr.Type)
	if err != nil {
		res.FailWithMsg("邮箱发送类型错误", c)
		return
	}
	switch emailType {
	case 1:
		// 查邮箱是否不存在
		var user models.UserModel
		err = global.DB.Take(&user, "email = ?", cr.Email).Error
		if err == nil {
			res.FailWithMsg("该邮箱已使用", c)
			return
		}
		err = email_service.SendRegisterCode(cr.Email, code)
	case 2:
		var user models.UserModel
		err = global.DB.Take(&user, "email = ?", cr.Email).Error
		if err != nil {
			res.FailWithMsg("该邮箱不存在", c)
			return
		}
		// 还必须得是邮箱注册的
		if user.RegisterSource != int8(enum.EmailLoginType) {
			res.FailWithMsg("非邮箱注册用户，不能重置密码", c)
			return
		}
		err = email_service.SendResetPasswordCode(cr.Email, code)
	case 3:
		var user models.UserModel
		err = global.DB.Take(&user, "email = ?", cr.Email).Error
		if err == nil {
			res.FailWithMsg("该邮箱已使用", c)
			return
		}
		err = email_service.SendBIndEmailCode(cr.Email, code)
	}
	if err != nil {
		slog.Error("邮件发送失败", "err", err)
		res.FailWithMsg("邮件发送失败", c)
		return
	}
	email_store.Set(id, cr.Email, code)
	res.OkWithData(SendEmailResponse{
		EmailID: id,
	}, c)
}
