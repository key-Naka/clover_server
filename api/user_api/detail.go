package user_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/utils/jwts"
	"time"

	"github.com/gin-gonic/gin"
)

type UserDetailResponse struct {
	ID             uint                    `json:"id"`
	CreatedAt      time.Time               `json:"createdAt"`
	Username       string                  `json:"username"`
	Nickname       string                  `json:"nickname"`
	Avatar         string                  `json:"avatar"`
	Abstract       string                  `json:"abstract"`
	RegisterSource enum.RegisterSourceType `json:"registerSource"` // 注册来源
	Role           enum.RoleType           `json:"role"`           // 角色
	models.UserConfModel
	Email       string `json:"email"`
	UsePassword bool   `json:"usePassword"`
}

func (UserApi) DetailView(c *gin.Context) {
	claims := jwts.GetClaims(c)
	var user models.UserModel
	err := global.DB.Preload("UserConf").Take(&user, claims.ID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}
	var data = UserDetailResponse{
		ID:             user.ID,
		CreatedAt:      user.CreatedAt,
		Username:       user.Username,
		Nickname:       user.Nickname,
		Avatar:         user.Avatar,
		Abstract:       user.Abstract,
		Role:           user.Role,
		RegisterSource: user.RegisterSource,
		Email:          user.Email,
	}

	if user.Password != "" {
		data.UsePassword = true
	}
	if user.UserConfModel != nil {
		data.UserConfModel = *user.UserConfModel
	}
	res.OkWithData(data, c)
}
