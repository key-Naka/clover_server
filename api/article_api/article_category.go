package article_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/utils/jwts"
	"strings"

	"github.com/gin-gonic/gin"
)

type CategoryCreateRequest struct {
	ID    uint   `json:"id"`
	Title string `json:"title" binding:"required"`
}

type CategoryListRequest struct {
	common.PageInfo
	UserID uint `form:"userID"`
	Type   int8 `form:"type" binding:"required,oneof=1 2 3"`
}

// CategoryCreateView 创建或更新当前用户的文章分类。
// @Summary 创建或更新文章分类
// @Description `id` 为 0 时创建分类；非 0 时更新当前用户已有分类。
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CategoryCreateRequest true "分类信息"
// @Success 200 {object} res.MessageResponse "分类创建或更新成功"
// @Failure 400 {object} res.ErrorResponse "参数、归属或认证失败"
// @Failure 500 {object} res.ErrorResponse "服务异常"
// @Router /article/category [post]
func (ArticleApi) CategoryCreateView(c *gin.Context) {
	var req CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	claims := jwts.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		res.FailWithMsg("分类标题不能为空", c)
		return
	}

	var category models.CategoryModel
	if req.ID != 0 {
		if err := global.DB.Take(&category, "id = ? and user_id = ?", req.ID, claims.UserID).Error; err != nil {
			res.FailWithMsg("分类不存在", c)
			return
		}
		if err := global.DB.Model(&category).Update("title", title).Error; err != nil {
			res.FailWithError(err, c)
			return
		}
		res.OkWithMsg("分类更新成功", c)
		return
	}

	if err := global.DB.Create(&models.CategoryModel{Title: title, UserID: claims.UserID}).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg("分类创建成功", c)
}

func (ArticleApi) CategoryListView(c *gin.Context) {
	var req CategoryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}

	query := global.DB.Model(&models.CategoryModel{}).Preload("User")
	switch req.Type {
	case 1:
		if req.UserID == 0 {
			res.FailWithMsg("用户id必填", c)
			return
		}
		query = query.Where("user_id = ?", req.UserID)
	case 2:
		claims := jwts.GetClaims(c)
		if claims == nil {
			res.FailWithMsg("请登录", c)
			return
		}
		query = query.Where("user_id = ?", claims.UserID)
	case 3:
		claims := jwts.GetClaims(c)
		if claims == nil || claims.Role != enum.AdminRole {
			res.FailWithMsg("角色错误", c)
			return
		}
		if req.UserID != 0 {
			query = query.Where("user_id = ?", req.UserID)
		}
	}
	if req.Key != "" {
		query = query.Where("title like ?", "%"+req.Key+"%")
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	var list []models.CategoryModel
	if err := query.Order("id desc").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&list).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithList(list, int(count), c)
}

func (ArticleApi) CategoryOptionsView(c *gin.Context) {
	claims := jwts.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}
	var list []models.CategoryModel
	if err := global.DB.Where("user_id = ?", claims.UserID).Order("id desc").Find(&list).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	options := make([]models.OptionsResponse[uint], 0, len(list))
	for _, item := range list {
		options = append(options, models.OptionsResponse[uint]{Label: item.Title, Value: item.ID})
	}
	res.OkWithData(options, c)
}
