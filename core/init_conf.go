package core

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	System System `yaml:"system"`
}

type System struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}


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
