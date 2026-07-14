// global/enter.go
package global

import (
	"clover_server/conf"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

const Version = "10.0.1"

var (
	Config *conf.Config
	Redis  *redis.Client
	DB     *gorm.DB
	DB1    *gorm.DB
)
