// conf/conf_system.go
package conf

// System 定义系统运行配置。
type System struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
}
