package focus_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

type FocusUserRequest struct {
	FocusUserID uint `json:"focusUserID" binding:"required"`
}

type FocusListRequest struct {
	common.PageInfo
	UserID uint `form:"userID"`
}

func (FocusApi) FocusUserView(c *gin.Context) {
	var req FocusUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	claims := jwts.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}
	if claims.UserID == req.FocusUserID {
		res.FailWithMsg("不能关注自己", c)
		return
	}
	var user models.UserModel
	if err := global.DB.Take(&user, req.FocusUserID).Error; err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}
	var focus models.UserFocusModel
	if err := global.DB.Take(&focus, "user_id = ? and focus_user_id = ?", claims.UserID, req.FocusUserID).Error; err == nil {
		res.FailWithMsg("已关注该用户", c)
		return
	}
	if err := global.DB.Create(&models.UserFocusModel{UserID: claims.UserID, FocusUserID: req.FocusUserID}).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg("关注成功", c)
}

func (FocusApi) UnFocusUserView(c *gin.Context) {
	var req FocusUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	claims := jwts.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}
	if claims.UserID == req.FocusUserID {
		res.FailWithMsg("不能取关自己", c)
		return
	}
	if err := global.DB.Where("user_id = ? and focus_user_id = ?", claims.UserID, req.FocusUserID).Delete(&models.UserFocusModel{}).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg("取消关注成功", c)
}

func (FocusApi) FocusUserListView(c *gin.Context) {
	var req FocusListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	claims, _ := jwts.ParseTokenByGin(c)
	userID := req.UserID
	if userID == 0 {
		if claims == nil {
			res.FailWithMsg("请登录", c)
			return
		}
		userID = claims.UserID
	}
	query := global.DB.Model(&models.UserFocusModel{}).Preload("FocusUser")
	query = query.Where("user_id = ?", userID)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	var list []models.UserFocusModel
	if err := query.Order("id desc").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&list).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithList(list, int(count), c)
}

func (FocusApi) FansUserListView(c *gin.Context) {
	var req FocusListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	claims, _ := jwts.ParseTokenByGin(c)
	userID := req.UserID
	if userID == 0 {
		if claims == nil {
			res.FailWithMsg("请登录", c)
			return
		}
		userID = claims.UserID
	}
	query := global.DB.Model(&models.UserFocusModel{}).Preload("User")
	query = query.Where("focus_user_id = ?", userID)
	var count int64
	if err := query.Count(&count).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	var list []models.UserFocusModel
	if err := query.Order("id desc").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&list).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithList(list, int(count), c)
}
