package core

import (
	"clover_server/conf"
	"clover_server/logger"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ToLoggerConfig 将业务配置转换为日志初始化配置。
func ToLoggerConfig(c *conf.Config) logger.Config {
	return logger.Config{
		Env:     c.System.Env,
		Level:   c.Log.Level,
		Format:  c.Log.Format,
		Dir:     c.Log.Dir,
		Prefix:  c.Log.Prefix,
		Console: c.Log.Console,
		File:    c.Log.File,
		Rotate: logger.RotateConfig{
			Enabled:    c.Log.Rotate.Enabled,
			Interval:   c.Log.Rotate.Interval,
			MaxAgeDays: c.Log.Rotate.MaxAgeDays,
		},
	}
}

// InitConf 读取并解析 YAML 配置文件。
func InitConf(path string) *conf.Config {
	confData, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取配置文件失败: %v\n", err)
		os.Exit(1)
	}

	var cfg conf.Config
	if err := yaml.Unmarshal(confData, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "解析yaml配置文件格式错误: %v\n", err)
		os.Exit(1)
	}

	return &cfg
}
