// global/enter.go
package global

import (
	"clover_server/conf"

	"gorm.io/gorm"
)

const Version = "10.0.1"

var (
	Config *conf.Config
	DB     *gorm.DB
	DB1    *gorm.DB
)
