package user_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/utils/jwts"
	"clover_server/utils/mps"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type UserInfoUpdateRequest struct {
	Username    *string   `json:"username" s-u:"username"`
	Nickname    *string   `json:"nickname" s-u:"nickname"`
	Avatar      *string   `json:"avatar" s-u:"avatar"`
	Abstract    *string   `json:"abstract" s-u:"abstract"`
	LikeTags    *[]string `json:"likeTags" s-u-c:"like_tags"`
	OpenCollect *bool     `json:"openCollect" s-u-c:"open_collect"`  // 公开我的收藏
	OpenFollow  *bool     `json:"openFollow" s-u-c:"open_follow"`    // 公开我的关注
	OpenFans    *bool     `json:"openFans" s-u-c:"open_fans"`        // 公开我的粉丝
	HomeStyleID *uint     `json:"homeStyleID" s-u-c:"home_style_id"` // 主页样式的id
}

func (UserApi) UpdateUserInfoView(c *gin.Context) {
	var req UserInfoUpdateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	userMap := mps.StructToMap(req, "s-u")
	userConfMap := mps.StructToMap(req, "s-u-c")
	claims := jwts.GetClaims(c)
	if len(userMap) > 0 {
		var userModel models.UserModel
		err := global.DB.Preload("UserConfModel").Take(&userModel, claims.UserID).Error
		if err != nil {
			res.FailWithError(err, c)
			return
		}
		if req.Username != nil {
			var UserCount int64
			global.DB.Debug().Model(models.UserModel{}).
				Where("username = ? and id <> ?", *req.Username, claims.UserID).
				Count(&UserCount)
			fmt.Println(*req.Username, UserCount)
			if UserCount > 0 {
				res.FailWithMsg("用户名已存在", c)
				return
			}
			if *req.Username != userModel.Username {
				// var uud = userModel.UserConfModel.UpdateUsernameDate
				// if uud != nil {
				// 	if time.Now().Sub(*uud).Hours() < 720 {
				// 		res.FailWithMsg("用户名30天内只能修改一次", c)
				// 		return
				// 	}
				// }
				userConfMap["update_username_date"] = time.Now()
			}
			if req.Nickname != nil || req.Avatar != nil {
				if userModel.RegisterSource == enum.RegisterQQSourceType {
					res.FailWithMsg("QQ注册的用户不能修改昵称和头像", c)
					return
				}
			}
		}
		err = global.DB.Model(&userModel).Updates(userMap).Error
		if err != nil {
			res.FailWithMsg("用户信息修改失败", c)
			return
		}
	}
	if len(userConfMap) > 0 {
		var userConfModel models.UserConfModel
		err = global.DB.Take(&userConfModel, "user_id = ?", claims.UserID).Error
		if err != nil {
			res.FailWithMsg("用户配置信息不存在", c)
			return
		}
		err = global.DB.Model(&userConfModel).Updates(userConfMap).Error
		if err != nil {
			res.FailWithMsg("用户信息修改失败", c)
			return
		}
	}
	res.OkWithMsg("用户信息修改成功", c)
}
