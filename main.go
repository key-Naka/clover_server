package main

import (
	"clover_server/core"
	"clover_server/flags"
	"clover_server/global"
	"clover_server/logger"

	"fmt"
	"log/slog"
	"os"
)

func main() {
	flags.ParseFlags()

	cnf, err := core.ReadConf(flags.FlagsOptions.File)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取配置失败: %v\n", err)
		os.Exit(1)
	}
	global.Config = cnf

	if err := logger.Init(core.ToLoggerConfig(cnf)); err != nil {
		fmt.Fprintf(os.Stderr, "初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := logger.Close(); err != nil {
			slog.Error("关闭日志资源失败", slog.Any("err", err))
		}
	}()
	slog.Info("成功加载配置",
		slog.String("config_file", flags.FlagsOptions.File),
		slog.String("env", cnf.System.Env),
		slog.String("ip", cnf.System.IP),
		slog.Int("port", cnf.System.Port),
		slog.String("log_level", cnf.Log.Level),
		slog.String("log_format", cnf.Log.Format),
	)
	global.DB = core.InitDB()
	slog.Info("数据库连接成功")
	// 1. 初始化 IP 数据库
	core.InitIPDB("112.45.173.110")

}
