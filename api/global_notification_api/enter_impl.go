package global_notification_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Href     string `json:"href"`
	Icon     string `json:"icon"`
}

type ListRequest struct {
	common.PageInfo
	Type int8 `form:"type" binding:"required,oneof=1 2"`
}

type UserActionRequest struct {
	ID   uint `json:"id" binding:"required"`
	Type int8 `json:"type" binding:"required,oneof=1 2"`
}

func (GlobalNotificationApi) CreateView(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	var count int64
	global.DB.Model(&models.GlobalNotificationModel{}).Where("title = ?", req.Title).Count(&count)
	if count > 0 {
		res.FailWithMsg("标题已存在", c)
		return
	}
	data := models.GlobalNotificationModel{Title: req.Title, Content: req.Content, Href: req.Href, Icon: req.Icon}
	if err := global.DB.Create(&data).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithData(data, c)
}

func (GlobalNotificationApi) ListView(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	claims := jwts.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}
	query := global.DB.Model(&models.GlobalNotificationModel{})
	if req.Type == 2 && claims.Role != enum.AdminRole {
		res.FailWithMsg("角色错误", c)
		return
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	var list []models.GlobalNotificationModel
	if err := query.Order("id desc").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&list).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithList(list, int(count), c)
}

func (GlobalNotificationApi) RemoveAdminView(c *gin.Context) {
	var req models.RemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	if err := global.DB.Delete(&models.GlobalNotificationModel{}, req.IDList).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg("删除成功", c)
}

func (GlobalNotificationApi) UserMsgActionView(c *gin.Context) {
	var req UserActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	claims := jwts.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}
	var notification models.GlobalNotificationModel
	if err := global.DB.Take(&notification, req.ID).Error; err != nil {
		res.FailWithMsg("通知不存在", c)
		return
	}
	var record models.UserGlobalNotificationModel
	_ = global.DB.Take(&record, "user_id = ? and notification_id = ?", claims.UserID, req.ID).Error
	record.UserID = claims.UserID
	record.NotificationID = req.ID
	if req.Type == 1 {
		record.IsRead = true
	}
	if req.Type == 2 {
		record.IsDelete = true
	}
	if record.ID == 0 {
		if err := global.DB.Create(&record).Error; err != nil {
			res.FailWithError(err, c)
			return
		}
	} else {
		if err := global.DB.Save(&record).Error; err != nil {
			res.FailWithError(err, c)
			return
		}
	}
	res.OkWithMsg("操作成功", c)
}
