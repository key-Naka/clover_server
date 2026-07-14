// conf/conf_jwt.go
package conf

type JWT struct {
	Expire int    `yaml:"expire"`
	Secret string `yaml:"secret"`
	Issuer string `yaml:"issuer"`
}
