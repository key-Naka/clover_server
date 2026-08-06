package user_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/utils/pwd"

	"github.com/gin-gonic/gin"
)

type ResetPasswordRequest struct {
	Pwd string `json:"pwd" binding:"required"`
}

// ResetPasswordSwaggerRequest 是携带邮箱验证码的重置密码请求。
type ResetPasswordSwaggerRequest struct {
	EmailID   string `json:"emailID" binding:"required"`
	EmailCode string `json:"emailCode" binding:"required"`
	Pwd       string `json:"pwd" binding:"required"`
}

// ResetPasswordView 通过邮箱验证码重置密码。
// @Summary 重置密码
// @Description 邮箱验证码中间件会从同一 JSON 请求体读取并校验 `emailID`、`emailCode`；校验成功后使用 `pwd` 重置该邮箱账户密码。
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body ResetPasswordSwaggerRequest true "邮箱验证码与新密码"
// @Success 200 {object} res.MessageResponse "重置密码成功"
// @Failure 400 {object} res.ErrorResponse "参数、邮箱验证码或账户状态校验失败"
// @Router /user/password/reset [put]
func (UserApi) ResetPasswordView(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}
	if !global.Config.Site.Login.EmailLogin {
		res.FailWithMsg("邮箱登录功能未开启", c)
		return
	}
	_email, _ := c.Get("email")
	email := _email.(string)
	var user models.UserModel
	err := global.DB.Take(&user, "email = ?", email).Error
	if err != nil {
		res.FailWithMsg("不存在的用户", c)
		return
	}
	if user.RegisterSource != enum.RegisterEmailSourceType {
		res.FailWithMsg("非邮箱注册用户，不能重置密码", c)
		return
	}
	hashPwd, _ := pwd.GenerateFromPassword(req.Pwd)
	global.DB.Model(&user).Update("password", hashPwd)
	res.OkWithMsg("重置密码成功", c)
}
