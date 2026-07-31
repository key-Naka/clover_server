package core

import (
	"clover_server/global"
	"log/slog"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

// InitES 初始化 Elasticsearch 客户端并执行连通性探测。
func InitES() *elasticsearch.Client {
	cfg := global.Config.Es
	if !cfg.Enable {
		slog.Info("elasticsearch未启用")
		return nil
	}

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.Addr()},
		Username:  cfg.Username,
		Password:  cfg.Password,
	})
	if err != nil {
		slog.Error("elasticsearch客户端初始化失败", slog.Any("err", err), slog.String("addr", cfg.Addr()))
		os.Exit(1)
	}

	resp, err := es.Info()
	if err != nil {
		slog.Error("elasticsearch连接失败", slog.Any("err", err), slog.String("addr", cfg.Addr()))
		os.Exit(1)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			slog.Warn("关闭elasticsearch响应体失败", slog.Any("err", closeErr))
		}
	}()

	if resp.IsError() {
		slog.Error("elasticsearch连通性检测失败", slog.String("addr", cfg.Addr()), slog.String("status", resp.Status()))
		os.Exit(1)
	}

	slog.Info("elasticsearch连接成功", slog.String("addr", cfg.Addr()), slog.String("status", resp.Status()))
	return es
}
