package core

import (
	"clover_server/logger"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	System System    `yaml:"system"`
	Log    LogConfig `yaml:"log"`
}

// System 定义系统运行配置。
type System struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
}

// LogConfig 定义日志输出与切分策略。
type LogConfig struct {
	Level   string       `yaml:"level"`
	Format  string       `yaml:"format"`
	Dir     string       `yaml:"dir"`
	Prefix  string       `yaml:"prefix"`
	Console bool         `yaml:"console"`
	File    bool         `yaml:"file"`
	Rotate  RotateConfig `yaml:"rotate"`
}

// RotateConfig 定义按时间滚动的日志策略。
type RotateConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Interval   string `yaml:"interval"`
	MaxAgeDays int    `yaml:"maxAgeDays"`
}

// ToLoggerConfig 将业务配置转换为日志初始化配置。
func (c *Config) ToLoggerConfig() logger.Config {
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

// ReadConf 读取并解析 YAML 配置文件。
func ReadConf(path string) (*Config, error) {
	confData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(confData, &cfg); err != nil {
		return nil, fmt.Errorf("解析yaml配置文件格式错误: %w", err)
	}

	return &cfg, nil
}
