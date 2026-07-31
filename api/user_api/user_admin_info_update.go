package user_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/utils/mps"

	"github.com/gin-gonic/gin"
)

type AdminUserInfoUpdateRequest struct {
	UserID   uint           `json:"userID" binding:"required"`
	Username *string        `json:"username" s-u:"username"`
	Nickname *string        `json:"nickname" s-u:"nickname"`
	Avatar   *string        `json:"avatar" s-u:"avatar"`
	Abstract *string        `json:"abstract" s-u:"abstract"`
	Role     *enum.RoleType `json:"role" s-u:"role"`
}

func (u *UserApi) AdminUserInfoUpdateView(c *gin.Context) {
	var req AdminUserInfoUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	userMap := mps.StructToMap(req, "s-u")
	var user models.UserModel
	err := global.DB.Take(&user, req.UserID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}
	err = global.DB.Model(&user).Updates(userMap).Error
	if err != nil {
		res.FailWithMsg("用户信息修改失败", c)
		return
	}
}
