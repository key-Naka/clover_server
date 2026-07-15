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
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (i ImageApi) ImageUploadView(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	uploadConf := global.Config.Upload
	maxSize := int64(uploadConf.Size) * 1024 * 1024
	if maxSize > 0 && fileHeader.Size > maxSize {
		res.FailWithMsg(fmt.Sprintf("文件大小不能超过%dMB", uploadConf.Size), c)
		return
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileHeader.Filename)), ".")
	if ext == "" || !isAllowedUploadExt(ext, uploadConf.WhiteList) {
		res.FailWithMsg(fmt.Sprintf("文件后缀必须为%s", strings.Join(uploadConf.WhiteList, "、")), c)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	fileHash := hash.Md5(fileBytes)
	var imageModel models.ImageModel
	err = global.DB.Where("hash = ?", fileHash).Take(&imageModel).Error
	if err == nil {
		res.OkWithData(BuildImageUploadResponse(imageModel), c)
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		res.FailWithError(err, c)
		return
	}

	saveFilename := fmt.Sprintf("%s.%s", fileHash, ext)
	savePath := filepath.Join(uploadConf.UploadDir, saveFilename)
	if err = c.SaveUploadedFile(fileHeader, savePath); err != nil {
		res.FailWithError(err, c)
		return
	}

	if global.Config.QiNiu.Enable {
		if _, err = qiniu_service.SendFile(savePath); err != nil {
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
		res.FailWithError(err, c)
		return
	}

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
