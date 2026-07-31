package conf

import "fmt"

type ESConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Scheme   string `yaml:"scheme"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Enable   bool   `yaml:"enable"`
}

func (c ESConfig) Addr() string {
	return fmt.Sprintf("%s://%s:%d", c.Scheme, c.Host, c.Port)
}
