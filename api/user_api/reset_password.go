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
