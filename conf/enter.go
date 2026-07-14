// conf/enter.go
package conf

type Config struct {
	System System      `yaml:"system"`
	Jwt    JWT         `yaml:"jwt"`
	Log    LogConfig   `yaml:"log"`
	Redis  RedisConfig `yaml:"redis"`
	DB     DBConfig    `yaml:"db"`
	DB1    DBConfig    `yaml:"db1"` // 数据库连接列表
}
