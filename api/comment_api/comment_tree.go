package comment_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

type CommentTreeItem struct {
	models.CommentModel
	Nickname string            `json:"nickname"`
	Avatar   string            `json:"avatar"`
	IsDigg   bool              `json:"isDigg"`
	Children []CommentTreeItem `json:"children"`
}

func (CommentApi) CommentTreeView(c *gin.Context) {
	var req models.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	article, err := loadArticle(req.ID)
	if err != nil || article.Status != 3 {
		res.FailWithMsg("文章不存在", c)
		return
	}
	var commentList []models.CommentModel
	if err = global.DB.Preload("UserModel").Where("article_id = ?", req.ID).Order("id asc").Find(&commentList).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	claims, _ := jwts.ParseTokenByGin(c)
	diggMap := map[uint]bool{}
	if claims != nil {
		var diggList []models.CommentDiggModel
		global.DB.Where("user_id = ?", claims.UserID).Find(&diggList)
		for _, item := range diggList {
			diggMap[item.CommentID] = true
		}
	}
	childrenMap := map[uint][]CommentTreeItem{}
	roots := make([]CommentTreeItem, 0)
	for _, item := range commentList {
		node := CommentTreeItem{
			CommentModel: item,
			Nickname:     item.UserModel.Nickname,
			Avatar:       item.UserModel.Avatar,
			IsDigg:       diggMap[item.ID],
		}
		if item.ParentID == nil {
			roots = append(roots, node)
			continue
		}
		childrenMap[*item.ParentID] = append(childrenMap[*item.ParentID], node)
	}
	var buildTree func(items []CommentTreeItem) []CommentTreeItem
	buildTree = func(items []CommentTreeItem) []CommentTreeItem {
		for index := range items {
			items[index].Children = buildTree(childrenMap[items[index].ID])
		}
		return items
	}
	res.OkWithData(buildTree(roots), c)
}
