package flags

import (
	"clover_server/global"
	"clover_server/models"
	"log/slog"
)

func FlagDB() {
	err := global.DB.AutoMigrate(&models.UserModel{})
	if err != nil {
		slog.Error("数据库迁移失败", "err", err)
		return
	}
	slog.Info("数据库迁移成功")
}
