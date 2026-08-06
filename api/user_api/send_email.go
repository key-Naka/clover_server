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

// SendEmail 发送邮箱验证码。
// @Summary 发送邮箱验证码
// @Description 请求体 `type`：1 注册、2 重置密码、3 绑定邮箱。当前路由已配置邮箱验证码中间件，该中间件要求请求中同时带有有效的 `emailID` 与 `emailCode`，因此直接首次调用会因中间件校验而无法进入本处理器；此处仅如实记录现状，未改变业务逻辑。
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body SendEmailRequest true "邮箱及验证码用途"
// @Success 200 {object} res.EmailCodeResponse "邮箱验证码标识"
// @Failure 400 {object} res.ErrorResponse "参数、邮箱状态或中间件验证码校验失败"
// @Failure 500 {object} res.ErrorResponse "邮件服务异常"
// @Router /user/send_email [post]
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
		if user.RegisterSource != enum.RegisterEmailSourceType {
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
