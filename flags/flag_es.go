// flags/flag_es.go
package flags

import (
	"clover_server/global"
	"clover_server/service/es_service"
	"log/slog"
)

func EsIndex() {
	if global.ES == nil {
		slog.Warn("未开启es连接")
		return
	}
	es_service.RebuildIndex(es_service.ArticleIndex, es_service.ArticleMapping)
	es_service.RebuildIndex(es_service.TextIndex, es_service.TextMapping)
}

func EsSync() {
	if global.ES == nil {
		slog.Warn("未开启es连接")
		return
	}
	if global.DB == nil {
		slog.Warn("未开启数据库连接")
		return
	}
	if err := es_service.SyncAllArticleDocuments(); err != nil {
		slog.Error("同步 MySQL 数据到 ES 失败", slog.String("error", err.Error()))
		return
	}
	slog.Info("同步 MySQL 数据到 ES 成功")
}
