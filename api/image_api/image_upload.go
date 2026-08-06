package image_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	qiniu_service "clover_server/service/qiniu_service"
	"clover_server/utils/hash"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ImageUploadView 上传图片。
// @Summary 上传图片
// @Description 仅管理员可调用。实际配置最大 10MB，仅接受 jpg、jpeg、png、webp、gif；服务以文件 MD5 去重，成功后返回本地及 Web 访问路径。
// @Tags 图片
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "图片文件"
// @Success 200 {object} res.ImageUploadResponse "上传成功"
// @Failure 400 {object} res.ErrorResponse "文件缺失、大小或格式不合法"
// @Failure 500 {object} res.ErrorResponse "存储服务异常"
// @Router /images [post]
func (i ImageApi) ImageUploadView(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		slog.Warn("获取上传文件失败", "path", c.Request.URL.Path, "client_ip", c.ClientIP(), "err", err)
		res.FailWithError(err, c)
		return
	}

	uploadConf := global.Config.Upload
	maxSize := int64(uploadConf.Size) * 1024 * 1024
	if maxSize > 0 && fileHeader.Size > maxSize {
		slog.Warn("上传文件大小超限", "filename", fileHeader.Filename, "size", fileHeader.Size, "max_size", maxSize)
		res.FailWithMsg(fmt.Sprintf("文件大小不能超过%dMB", uploadConf.Size), c)
		return
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileHeader.Filename)), ".")
	if ext == "" || !isAllowedUploadExt(ext, uploadConf.WhiteList) {
		slog.Warn("上传文件后缀不允许", "filename", fileHeader.Filename, "ext", ext, "allow_list", uploadConf.WhiteList)
		res.FailWithMsg(fmt.Sprintf("文件后缀必须为%s", strings.Join(uploadConf.WhiteList, "、")), c)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		slog.Error("打开上传文件失败", "filename", fileHeader.Filename, "err", err)
		res.FailWithError(err, c)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		slog.Error("读取上传文件失败", "filename", fileHeader.Filename, "err", err)
		res.FailWithError(err, c)
		return
	}

	fileHash := hash.Md5(fileBytes)
	var imageModel models.ImageModel
	err = global.DB.Where("hash = ?", fileHash).Take(&imageModel).Error
	if err == nil {
		slog.Info("图片重复上传，返回已存在文件", "filename", fileHeader.Filename, "hash", fileHash, "image_id", imageModel.ID)
		res.OkWithData(BuildImageUploadResponse(imageModel), c)
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Error("查询图片哈希失败", "filename", fileHeader.Filename, "hash", fileHash, "err", err)
		res.FailWithError(err, c)
		return
	}

	saveFilename := fmt.Sprintf("%s.%s", fileHash, ext)
	savePath := filepath.Join(uploadConf.UploadDir, saveFilename)
	if err = c.SaveUploadedFile(fileHeader, savePath); err != nil {
		slog.Error("保存上传文件失败", "filename", fileHeader.Filename, "save_path", savePath, "err", err)
		res.FailWithError(err, c)
		return
	}

	if global.Config.QiNiu.Enable {
		if _, err = qiniu_service.SendFile(savePath); err != nil {
			slog.Error("上传文件到七牛失败", "filename", saveFilename, "save_path", savePath, "err", err)
			res.FailWithError(err, c)
			return
		}
	}

	imageModel = models.ImageModel{
		Filename: saveFilename,
		Path:     savePath,
		Size:     fileHeader.Size,
		Hash:     fileHash,
	}
	if err = global.DB.Create(&imageModel).Error; err != nil {
		slog.Error("保存图片记录失败", "filename", saveFilename, "path", savePath, "hash", fileHash, "err", err)
		res.FailWithError(err, c)
		return
	}

	slog.Info("图片上传成功", "image_id", imageModel.ID, "filename", saveFilename, "path", savePath, "size", fileHeader.Size, "hash", fileHash)
	res.OkWithData(BuildImageUploadResponse(imageModel), c)
}

func isAllowedUploadExt(ext string, whiteList []string) bool {
	for _, allowExt := range whiteList {
		if ext == strings.TrimPrefix(strings.ToLower(allowExt), ".") {
			return true
		}
	}
	return false
}

func BuildImageUploadResponse(imageModel models.ImageModel) gin.H {
	return gin.H{
		"path":     imageModel.Path,
		"web_path": imageModel.WebPath(),
		"filename": imageModel.Filename,
		"size":     imageModel.Size,
		"hash":     imageModel.Hash,
	}
}
