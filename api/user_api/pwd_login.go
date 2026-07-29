package user_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/service/log_service"
	"clover_server/service/user_service"
	"clover_server/utils/jwts"
	"clover_server/utils/pwd"

	"github.com/gin-gonic/gin"
)

type PwdLoginRequest struct {
	Val      string `json:"val" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (UserApi) PwdLoginApi(c *gin.Context) {
	var cr PwdLoginRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	if !global.Config.Site.Login.UsernamePwdLogin {
		res.FailWithMsg("站点未启用密码登陆", c)
		return
	}
	var user models.UserModel
	err = global.DB.Take(&user, "(username = ? or email = ?) and password  <> ''", cr.Val, cr.Val).Error
	if err != nil {
		log_service.NewLoginFail(c, enum.UserPwdLoginType, "用户不存在", cr.Val, cr.Password)
		res.FailWithMsg("用户不存在或密码错误", c)
		return
	}
	if !pwd.CompareHashAndPassword(user.Password, cr.Password) {
		log_service.NewLoginFail(c, enum.UserPwdLoginType, "密码错误", cr.Val, cr.Password)
		res.FailWithMsg("用户不存在或密码错误", c)
		return
	}
	token, _ := jwts.GetToken(jwts.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     enum.RoleType(user.Role),
	})
	user_service.NewUserService(user).UserLogin(c)

	res.OkWithData(token, c)
}
