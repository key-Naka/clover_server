package image_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/service/log_service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ImageApi struct {
}

type ImageListViewResponse struct {
	ImageModel models.ImageModel
	WebPath    string `json:"web_path"`
}

func (ImageApi) ImageListView(c *gin.Context) {
	var pageInfo common.PageInfo
	c.ShouldBindQuery(&pageInfo)
	list, count, err := common.ListQuery(models.ImageModel{}, common.Options{
		PageInfo: pageInfo,
		Likes:    []string{"filename"},
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	var listResponse = make([]ImageListViewResponse, 0)
	for _, item := range list {
		listResponse = append(listResponse, ImageListViewResponse{
			ImageModel: item,
			WebPath:    item.WebPath(),
		})
	}

	res.OkWithList(listResponse, count, c)
}

func (ImageApi) ImageRemoveView(c *gin.Context) {
	var req models.RemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	if len(req.IDList) == 0 {
		res.FailWithMsg("请选择要删除的图片", c)
		return
	}

	log := log_service.GetLog(c)
	log.ShowRequest()
	log.ShowResponse()

	var imageModel []models.ImageModel
	result := global.DB.Find(&imageModel, req.IDList)
	if result.Error != nil {
		res.FailWithMsg("查询图片失败", c)
		return
	}
	if len(imageModel) == 0 {
		res.FailWithMsg("图片不存在", c)
		return
	}
	if int64(len(imageModel)) != int64(len(req.IDList)) {
		res.FailWithMsg("部分图片不存在，删除失败", c)
		return
	}

	if err := global.DB.Delete(&imageModel).Error; err != nil {
		res.FailWithMsg("删除图片失败", c)
		return
	}

	msg := fmt.Sprintf("删除%d条图片成功", len(imageModel))
	res.OkWithMsg(msg, c)
}
