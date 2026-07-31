// global/enter.go
package global

import (
	"clover_server/conf"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"github.com/mojocn/base64Captcha"
	"gorm.io/gorm"
)

const Version = "10.0.1"

var (
	Config       *conf.Config
	Qiniu        *conf.QiNiu
	Redis        *redis.Client
	ES           *elasticsearch.Client
	DB           *gorm.DB
	DB1          *gorm.DB
	CaptchaStore = base64Captcha.DefaultMemStore
)
