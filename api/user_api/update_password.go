package user_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/utils/jwts"
	"clover_server/utils/pwd"

	"github.com/gin-gonic/gin"
)

type UpdatePasswordRequest struct {
	OldPwd string `json:"oldPwd" binding:"required"`
	Pwd    string `json:"pwd" binding:"required"`
}

func (UserApi) UpdatePasswordView(c *gin.Context) {
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}
	claims := jwts.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("用户不存在", c)

		return
	}
	var user models.UserModel
	err := global.DB.Take(&user, "id = ?", claims.ID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}
	if !pwd.CompareHashAndPassword(user.Password, req.OldPwd) {
		res.FailWithMsg("旧密码错误", c)
		return
	}
	hashPwd, _ := pwd.GenerateFromPassword(req.Pwd)
	global.DB.Model(&user).Update("password", hashPwd)
	res.OkWithMsg("密码修改成功", c)
}
