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
	es_service.RebuildIndex("article_index", es_service.ArticleMapping)
	es_service.RebuildIndex("text_index", es_service.TextMapping)
}
