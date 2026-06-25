package conf

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
