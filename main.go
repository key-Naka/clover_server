package main

import (
	"clover_server/core"
	"clover_server/flags"
	"clover_server/global"
	"clover_server/logger"
	"clover_server/router"
	"log/slog"
)

// @title Clover Server API
// @version 1.0
// @description Clover Server 的 REST API 文档。
// @BasePath /api
// @schemes https http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入 Bearer {JWT}，例如：Bearer eyJhbGciOiJIUzI1NiIs...
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
	global.Redis = core.InitRedis()
	global.ES = core.InitES()
	core.InitIPDB()
	flags.Run()

	router.Run()
}
