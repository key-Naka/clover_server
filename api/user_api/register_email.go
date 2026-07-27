// api/user_api/register_email.go
package user_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/service/user_service"
	"clover_server/utils/jwts"
	"clover_server/utils/pwd"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type RegisterEmailRequest struct {
	EmailID   string `json:"emailID" binding:"required"`
	EmailCode string `json:"emailCode" binding:"required"`
	Pwd       string `json:"pwd" binding:"required"`
}

func (UserApi) RegisterEmailView(c *gin.Context) {
	var cr RegisterEmailRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	if !global.Config.Site.Login.EmailLogin {
		res.FailWithMsg("站点未启用邮箱注册", c)
		return
	}
	uname := base64Captcha.RandText(5, "0123456789")

	_email, _ := c.Get("email")
	email := _email.(string)

	hashPwd, _ := pwd.GenerateFromPassword(cr.Pwd)
	var user = models.UserModel{
		Username:       fmt.Sprintf("b_%s", uname),
		Nickname:       "邮箱用户",
		RegisterSource: int8(enum.EmailLoginType),
		Password:       hashPwd,
		Email:          email,
		Role:           int8(enum.UserRole),
	}

	err = global.DB.Create(&user).Error
	if err != nil {
		res.FailWithMsg("邮箱注册失败", c)
		slog.Error("创建用户失败", "err", err)
		return
	}

	token, err := jwts.GetToken(jwts.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     enum.RoleType(user.Role),
	})
	if err != nil {
		res.FailWithMsg("邮箱登录失败", c)
		return
	}
	user_service.NewUserService(user).UserLogin(c)
	res.OkWithData(token, c)
}
