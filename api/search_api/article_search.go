package search_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/service/es_service"

	"github.com/gin-gonic/gin"
)

type ArticleSearchRequest struct {
	common.PageInfo
	CategoryID *uint `form:"categoryID"`
}

func (SearchApi) ArticleSearchView(c *gin.Context) {
	var req ArticleSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	status := int8(3)
	response, err := es_service.SearchArticles(es_service.SearchArticlesRequest{
		Keyword:    req.Key,
		CategoryID: req.CategoryID,
		Status:     &status,
		Page:       req.Page,
		Limit:      req.Limit,
	})
	if err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}
	res.OkWithData(response, c)
}
