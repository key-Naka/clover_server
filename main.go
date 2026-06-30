package main

import (
	"clover_server/core"
	"clover_server/flags"
	"clover_server/global"
	"clover_server/logger"

	"log/slog"
)

func main() {
	flags.ParseFlags()

	cnf := core.InitConf(flags.FlagsOptions.File)
	global.Config = cnf

	logger.Init(core.ToLoggerConfig(cnf))
	defer func() {
		if err := logger.Close(); err != nil {
			slog.Error("关闭日志资源失败", slog.Any("err", err))
		}
	}()

	global.DB = core.InitDB()

	flags.Run()
}
