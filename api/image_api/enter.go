package image_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/service/log_service"
	"fmt"
	"log/slog"

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
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		slog.Warn("绑定图片列表查询参数失败", "path", c.Request.URL.Path, "client_ip", c.ClientIP(), "err", err)
	}
	list, count, err := common.ListQuery(models.ImageModel{}, common.Options{
		PageInfo: pageInfo,
		Likes:    []string{"filename"},
	})
	if err != nil {
		slog.Error("获取图片列表失败", "err", err, "page", pageInfo.Page, "limit", pageInfo.Limit, "key", pageInfo.Key)
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

// ImageRemoveView 批量删除图片。
// @Summary 删除图片
// @Tags 图片
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RemoveRequest true "待删除的图片 ID 列表"
// @Success 200 {object} res.MessageResponse "图片删除成功"
// @Failure 400 {object} res.ErrorResponse "参数、认证或图片存在性校验失败"
// @Failure 500 {object} res.ErrorResponse "删除失败"
// @Router /images [delete]
func (ImageApi) ImageRemoveView(c *gin.Context) {
	var req models.RemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("绑定删除图片参数失败", "path", c.Request.URL.Path, "client_ip", c.ClientIP(), "err", err)
		res.FailWithError(err, c)
		return
	}
	if len(req.IDList) == 0 {
		slog.Warn("删除图片失败，未选择图片", "path", c.Request.URL.Path, "client_ip", c.ClientIP())
		res.FailWithMsg("请选择要删除的图片", c)
		return
	}

	log := log_service.GetLog(c)
	log.ShowRequest()
	log.ShowResponse()
	log.SetTitle("删除图片")
	log.SetItemInfo("图片ID列表", req.IDList)

	var imageModel []models.ImageModel
	result := global.DB.Find(&imageModel, req.IDList)
	if result.Error != nil {
		slog.Error("查询待删除图片失败", "err", result.Error, "id_list", req.IDList)
		res.FailWithMsg("查询图片失败", c)
		return
	}
	if len(imageModel) == 0 {
		slog.Warn("删除图片失败，图片不存在", "id_list", req.IDList)
		res.FailWithMsg("图片不存在", c)
		return
	}
	if int64(len(imageModel)) != int64(len(req.IDList)) {
		slog.Warn("删除图片失败，存在部分图片缺失", "request_count", len(req.IDList), "found_count", len(imageModel), "id_list", req.IDList)
		res.FailWithMsg("部分图片不存在，删除失败", c)
		return
	}

	if err := global.DB.Delete(&imageModel).Error; err != nil {
		slog.Error("删除图片失败", "err", err, "id_list", req.IDList)
		res.FailWithMsg("删除图片失败", c)
		return
	}

	log.SetItemInfo("删除数量", len(imageModel))
	slog.Info("删除图片成功", "count", len(imageModel), "id_list", req.IDList)
	msg := fmt.Sprintf("删除%d条图片成功", len(imageModel))
	res.OkWithMsg(msg, c)
}
