package redis_jwt

import (
	"clover_server/global"
	"clover_server/utils/jwts"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

type BlackType int8

const (
	UserBlackType   BlackType = 1 // 用户注销登陆
	AdminBlackType  BlackType = 2 // 管理员让你手动下线
	DeviceBlackType BlackType = 3 // 其他设备把自己挤下来了
)

func (b BlackType) String() string {
	return fmt.Sprintf("%d", b)
}

func (b BlackType) Msg() string {
	switch b {
	case UserBlackType:
		return "已注销"
	case AdminBlackType:
		return "禁止登录"
	case DeviceBlackType:
		return "设备下线"
	}
	return "已注销"
}
func ParseBlackType(val string) BlackType {
	switch val {
	case "1":
		return UserBlackType
	case "2":
		return AdminBlackType
	case "3":
		return DeviceBlackType
	}
	return UserBlackType
}

func TokenBlack(token string, value BlackType) {
	key := fmt.Sprintf("token_black_%s", token)

	claims, err := jwts.ParseToken(token)
	if err != nil || claims == nil {
		slog.Error("token解析失败，加入黑名单终止", "err", err)
		return
	}

	second := claims.ExpiresAt.Unix() - time.Now().Unix()
	if second <= 0 {
		slog.Warn("token已过期，无需加入黑名单", "user_id", claims.UserID, "username", claims.Username, "black_type", value.String())
		return
	}

	_, err = global.Redis.Set(key, value.String(), time.Duration(second)*time.Second).Result()
	if err != nil {
		slog.Error("redis添加黑名单失败", "err", err, "user_id", claims.UserID, "username", claims.Username, "black_type", value.String())
		return
	}

	slog.Info("token加入黑名单成功", "user_id", claims.UserID, "username", claims.Username, "black_type", value.String(), "expired_in_seconds", second)
}

func HasTokenBlack(token string) (blk BlackType, ok bool) {
	key := fmt.Sprintf("token_black_%s", token)
	value, err := global.Redis.Get(key).Result()
	if err != nil {
		return
	}
	blk = ParseBlackType(value)
	return blk, true
}

func HasTokenBlackByGin(c *gin.Context) (blk BlackType, ok bool) {
	token := c.GetHeader("token")
	if token == "" {
		token = c.Query("token")
	}
	return HasTokenBlack(token)
}
