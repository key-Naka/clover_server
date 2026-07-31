package flags

import (
	"clover_server/global"
	"clover_server/models"
	"log/slog"

	"gorm.io/plugin/dbresolver"
)

func FlagDB() {
	err := global.DB.Clauses(dbresolver.Write).AutoMigrate(&models.UserModel{},
		&models.UserConfModel{},
		&models.ArticleModel{},
		&models.ArticleDiggModel{},
		&models.UserArticleLookHistoryModel{},
		&models.CategoryModel{},
		&models.ImageModel{},
		&models.UserArticleCollectModel{},
		&models.CollectModel{},
		&models.CommentModel{},
		&models.BannerModel{},
		&models.LogModel{},
		&models.GlobalNotificationModel{},
		&models.LogModel{},
		&models.UserLoginModel{},
	)
	if err != nil {
		slog.Error("数据库迁移失败", "err", err)
		return
	} else {
		slog.Info("数据库迁移成功")
	}
}
