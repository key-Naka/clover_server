package comment_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

type CommentListRequest struct {
	common.PageInfo
	ArticleID uint `form:"articleID"`
	UserID    uint `form:"userID"`
	Type      int8 `form:"type" binding:"required,oneof=1 2 3"`
}

type CommentListResponse struct {
	models.CommentModel
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	ArticleTitle string `json:"articleTitle"`
}

func (CommentApi) CommentListView(c *gin.Context) {
	var req CommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	claims := jwts.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}
	query := global.DB.Model(&models.CommentModel{}).Preload("UserModel")
	switch req.Type {
	case 1:
		var articleIDs []uint
		global.DB.Model(&models.ArticleModel{}).Where("user_id = ?", claims.UserID).Pluck("id", &articleIDs)
		query = query.Where("article_id in ?", articleIDs)
	case 2:
		query = query.Where("user_id = ?", claims.UserID)
	case 3:
		if claims.Role != enum.AdminRole {
			res.FailWithMsg("角色错误", c)
			return
		}
		if req.UserID != 0 {
			query = query.Where("user_id = ?", req.UserID)
		}
	}
	if req.ArticleID != 0 {
		query = query.Where("article_id = ?", req.ArticleID)
	}
	if req.Key != "" {
		query = query.Where("content like ?", "%"+req.Key+"%")
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	var list []models.CommentModel
	if err := query.Order("id desc").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&list).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	responseList := make([]CommentListResponse, 0, len(list))
	for _, item := range list {
		responseList = append(responseList, CommentListResponse{
			CommentModel: item,
			Nickname:     item.UserModel.Nickname,
			Avatar:       item.UserModel.Avatar,
			ArticleTitle: queryArticleTitle(item.ArticleID),
		})
	}
	res.OkWithList(responseList, int(count), c)
}

func queryArticleTitle(articleID uint) string {
	var article models.ArticleModel
	if err := global.DB.Take(&article, articleID).Error; err != nil {
		return ""
	}
	return article.Title
}
