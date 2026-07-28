package user_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/service/qq_service"
	"clover_server/service/user_service"
	"clover_server/utils/jwts"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type QQLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

func (u *UserApi) QQLoginView(c *gin.Context) {
	var cr QQLoginRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	if !global.Config.Site.Login.QQLogin {
		res.FailWithError(errors.New("qq登录未开启"), c)
		return
	}
	info, err := qq_service.GetUserInfo(cr.Code)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	var user models.UserModel
	err = global.DB.Take(&user, "open_id = ?", info.OpenID).Error
	if err != nil {
		uname := base64Captcha.RandText(5, "0123456789")
		user = models.UserModel{
			Username:       fmt.Sprintf("c_%s", uname),
			Nickname:       info.Nickname,
			Avatar:         info.Avatar,
			RegisterSource: int8(enum.RegisterQQSourceType),
			OpenID:         info.OpenID,
			Role:           int8(enum.UserRole),
		}
		err = global.DB.Create(&user).Error
		if err != nil {
			res.FailWithMsg("qq登录失败", c)
			return
		}
		token, _ := jwts.GetToken(jwts.Claims{
			Username: user.Username,
			Role:     enum.RoleType(user.Role),
			UserID:   user.ID,
		})
		user_service.NewUserService(user).UserLogin(c)
		res.OkWithData(token, c)
	}
}
