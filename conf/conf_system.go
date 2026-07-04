// conf/conf_system.go
package conf

import "fmt"

// System 定义系统运行配置。
type System struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
}

func (s *System) GetAddr() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}
