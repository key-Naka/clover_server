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

// UpdatePasswordView 修改当前登录用户密码。
// @Summary 修改密码
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdatePasswordRequest true "旧密码与新密码"
// @Success 200 {object} res.MessageResponse "密码修改成功"
// @Failure 400 {object} res.ErrorResponse "参数、认证或旧密码校验失败"
// @Router /user/password [put]
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
