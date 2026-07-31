// service/es_service/index.go
package es_service

import (
	"clover_server/global"
	"context"
	"log/slog"
	"strings"
)

// RebuildIndex 重建索引，不存在时直接创建，存在时先删除再创建。
func RebuildIndex(index, mapping string) {
	if ExistsIndex(index) {
		DeleteIndex(index)
	}
	CreateIndex(index, mapping)
}

// CreateIndex 创建 ES 索引。
func CreateIndex(index, mapping string) {
	_, err := global.ES.
		Indices.Create(index, global.ES.Indices.Create.WithContext(context.Background()), global.ES.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		slog.Error("索引创建失败", slog.String("index", index), slog.String("error", err.Error()))
		return
	}
	slog.Info("索引创建成功", slog.String("index", index))
}

// ExistsIndex 判断索引是否存在。
func ExistsIndex(index string) bool {
	response, err := global.ES.Indices.Exists([]string{index}, global.ES.Indices.Exists.WithContext(context.Background()))
	if err != nil {
		return false
	}
	defer response.Body.Close()
	return response.StatusCode == 200
}

// DeleteIndex 删除 ES 索引。
func DeleteIndex(index string) {
	response, err := global.ES.Indices.Delete([]string{index}, global.ES.Indices.Delete.WithContext(context.Background()))
	if err != nil {
		slog.Error("索引删除失败", slog.String("index", index), slog.String("error", err.Error()))
		return
	}
	defer response.Body.Close()
	slog.Info("索引删除成功", slog.String("index", index))
}
