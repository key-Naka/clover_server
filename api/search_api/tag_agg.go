package search_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"sort"

	"github.com/gin-gonic/gin"
)

type TagAggResponse struct {
	Tag          string `json:"tag"`
	ArticleCount int    `json:"articleCount"`
}

func (SearchApi) TagAggView(c *gin.Context) {
	var req common.PageInfo
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	var articleList []models.ArticleModel
	if err := global.DB.Where("status = ?", 3).Find(&articleList).Error; err != nil {
		res.FailWithError(err, c)
		return
	}
	tagMap := map[string]int{}
	for _, article := range articleList {
		for _, tag := range article.TagList {
			if tag == "" {
				continue
			}
			tagMap[tag]++
		}
	}
	list := make([]TagAggResponse, 0, len(tagMap))
	for tag, count := range tagMap {
		list = append(list, TagAggResponse{Tag: tag, ArticleCount: count})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].ArticleCount == list[j].ArticleCount {
			return list[i].Tag < list[j].Tag
		}
		return list[i].ArticleCount > list[j].ArticleCount
	})
	offset := req.GetOffset()
	if offset >= len(list) {
		res.OkWithList([]TagAggResponse{}, len(list), c)
		return
	}
	end := offset + req.GetLimit()
	if end > len(list) {
		end = len(list)
	}
	res.OkWithList(list[offset:end], len(list), c)
}
