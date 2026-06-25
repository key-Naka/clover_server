// conf/enter.go
package conf

type Config struct {
	System System    `yaml:"system"`
	Log    LogConfig `yaml:"log"`
	DB     DBConfig  `yaml:"db"`
	DB1    DBConfig  `yaml:"db1"` // 数据库连接列表
}
