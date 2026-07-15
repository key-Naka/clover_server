// core/init_redis.go
package core

import (
	"clover_server/global"

	"log/slog"

	"github.com/go-redis/redis"
)

func InitRedis() *redis.Client {
	r := global.Config.Redis
	redisDB := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})
	_, err := redisDB.Ping().Result()
	if err != nil {
		slog.Error("redis连接失败", "err", err)
		return nil
	}
	slog.Info("redis连接成功")
	return redisDB
}
