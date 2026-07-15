package models

import (
	"fmt"
	"os"
	"strings"

	"clover_server/global"
	"log/slog"

	"gorm.io/gorm"
)

type ImageModel struct {
	Model
	Filename string `gorm:"size:64" json:"filename"`
	Path     string `gorm:"size:256" json:"path"`
	Size     int64  `json:"size"`
	Hash     string `gorm:"size:32;index" json:"hash"`
}

func (i ImageModel) WebPath() string {
	if global.Config != nil && global.Config.QiNiu.Enable {
		baseURL := strings.TrimRight(global.Config.QiNiu.Uri, "/")
		prefix := strings.Trim(global.Config.QiNiu.Prefix, "/")
		if prefix == "" {
			return fmt.Sprintf("%s/%s", baseURL, i.Filename)
		}
		return fmt.Sprintf("%s/%s/%s", baseURL, prefix, i.Filename)
	}
	return fmt.Sprintf("/%s", i.Path)
}

func (l ImageModel) BeforeDelete(tx *gorm.DB) error {
	err := os.Remove(l.Path)
	if err != nil {
		slog.Error("删除文件失败", "path", l.Path, "error", err)
	}
	return nil
}
