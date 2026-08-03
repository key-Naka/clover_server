package search_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/service/es_service"

	"github.com/gin-gonic/gin"
)

type TextSearchRequest struct {
	common.PageInfo
}

func (SearchApi) TextSearchView(c *gin.Context) {
	var req TextSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	response, err := es_service.SearchTexts(es_service.SearchTextsRequest{Keyword: req.Key, Page: req.Page, Limit: req.Limit})
	if err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}
	res.OkWithData(response, c)
}
