package conf

import "fmt"

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
	Debug    bool   `yaml:"debug"`
	Source   string `yaml:"source"`
}

func (dbconfig DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		dbconfig.Username, dbconfig.Password, dbconfig.Host, dbconfig.Port, dbconfig.DBName, dbconfig.Charset)
}
