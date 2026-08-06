package article_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

type ArticleListRequest struct {
	common.PageInfo
	Type       int8  `form:"type" binding:"required,oneof=1 2 3"`
	UserID     uint  `form:"userID"`
	CategoryID *uint `form:"categoryID"`
	Status     *int8 `form:"status"`
}

type ArticleListResponse struct {
	models.ArticleModel
	CategoryTitle string `json:"categoryTitle,omitempty"`
	UserNickname  string `json:"userNickname"`
	UserAvatar    string `json:"userAvatar"`
}

// ArticleListView 获取文章列表。
// @Summary 获取文章列表
// @Description `type=1` 为指定用户已发布文章，必须提供 `userID`；`type=2` 为当前登录用户文章，需 Bearer JWT；`type=3` 为管理员视图，需管理员 Bearer JWT。`status` 用于状态筛选。
// @Tags 文章
// @Produce json
// @Param page query int false "页码，默认 1"
// @Param limit query int false "每页数量，默认 10，最大 100"
// @Param key query string false "标题或摘要关键字"
// @Param order query string false "排序字段，默认 id desc"
// @Param type query int true "列表类型：1 用户公开文章；2 我的文章；3 管理员文章"
// @Param userID query int false "用户 ID；type=1 时必填，type=3 时可选"
// @Param categoryID query int false "分类 ID"
// @Param status query int false "文章状态"
// @Success 200 {object} res.ArticleListResponse "data 包含 list 与 count"
// @Failure 400 {object} res.ErrorResponse "参数或权限校验失败"
// @Failure 500 {object} res.ErrorResponse "服务异常"
// @Router /article/list [get]
func (ArticleApi) ArticleListView(c *gin.Context) {
	var req ArticleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}

	query := global.DB.Model(&models.ArticleModel{}).Preload("UserModel")
	switch req.Type {
	case 1:
		if req.UserID == 0 {
			res.FailWithMsg("用户id必填", c)
			return
		}
		query = query.Where("user_id = ?", req.UserID).Where("status = ?", 3)
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

	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.Key != "" {
		like := "%" + req.Key + "%"
		query = query.Where("title like ? or abstract like ?", like, like)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		res.FailWithError(err, c)
		return
	}

	var list []models.ArticleModel
	order := req.Order
	if order == "" {
		order = "id desc"
	}
	if err := query.Order(order).Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&list).Error; err != nil {
		res.FailWithError(err, c)
		return
	}

	responseList := make([]ArticleListResponse, 0, len(list))
	for _, article := range list {
		item := ArticleListResponse{
			ArticleModel:  article,
			UserNickname:  article.UserModel.Nickname,
			UserAvatar:    article.UserModel.Avatar,
			CategoryTitle: queryCategoryTitle(article.CategoryID),
		}
		item.Content = ""
		responseList = append(responseList, item)
	}
	res.OkWithList(responseList, int(count), c)
}
