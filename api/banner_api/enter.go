package banner_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

type BannerApi struct {
}
type BannerCreateRequest struct {
	Title string `json:"title"`
	Cover string `json:"cover"`
	Href  string `json:"href"`
	Show  bool   `json:"show"`
}

func (i BannerApi) BannerCreate(c *gin.Context) {
	var req BannerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	err := global.DB.Create(&models.BannerModel{
		Title: req.Title,
		Cover: req.Cover,
		Href:  req.Href,
		Show:  req.Show,
	}).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithData(req, c)
}

func (i BannerApi) BannerList(c *gin.Context) {
	type bannerListRequest struct {
		common.PageInfo
		Show bool `json:"show"`
	}
	var cr bannerListRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(err, c)
		return
	}
	list, count, err := common.ListQuery(models.BannerModel{Show: cr.Show}, common.Options{
		PageInfo: cr.PageInfo,
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithList(list, count, c)
}

type BannerRemoveRequest struct {
	Ids []uint `json:"ids"`
}

func (i BannerApi) BannerRemove(c *gin.Context) {
	var req BannerRemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}

	if len(req.Ids) == 0 {
		res.FailWithError(fmt.Errorf("ids 不能为空"), c)
		return
	}

	result := global.DB.Delete(&models.BannerModel{}, req.Ids)
	if result.Error != nil {
		res.FailWithError(result.Error, c)
		return
	}

	res.OkWithData(fmt.Sprintf("删除成功%d个", result.RowsAffected), c)
}
func (i BannerApi) BannerUpdate(c *gin.Context) {
	var id models.IDRequest
	if err := c.ShouldBindUri(&id); err != nil {
		res.FailWithError(err, c)
		return
	}
	var banner models.BannerModel
	err := global.DB.Take(&banner, id.ID).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	var cr BannerCreateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	updateData := map[string]any{
		"title": cr.Title,
		"cover": cr.Cover,
		"href":  cr.Href,
		"show":  cr.Show,
	}
	err = global.DB.Model(&banner).Updates(updateData).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg("更新成功", c)
}
