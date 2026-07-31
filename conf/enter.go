// conf/enter.go
package conf

type Config struct {
	System System      `yaml:"system"`
	Jwt    JWT         `yaml:"jwt"`
	Log    LogConfig   `yaml:"log"`
	Redis  RedisConfig `yaml:"redis"`
	Es     ESConfig    `yaml:"es"`
	DB     DBConfig    `yaml:"db"`
	DB1    DBConfig    `yaml:"db1"` // 数据库连接列表
	Site   Site        `yaml:"site"`
	Email  Email       `yaml:"email"`
	QQ     QQ          `yaml:"qq"`
	QiNiu  QiNiu       `yaml:"qiNiu"`
	Ai     Ai          `yaml:"ai"`
	Upload Upload      `yaml:"upload"`
}
