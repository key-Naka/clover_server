package middleware

import (
	"clover_server/common/res"
	"clover_server/models/enum"
	"clover_server/service/redis_service/redis_jwt"
	"clover_server/utils/jwts"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		slog.Warn("用户鉴权失败", "path", c.Request.URL.Path, "method", c.Request.Method, "client_ip", c.ClientIP(), "err", err)
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	blk, ok := redis_jwt.HasTokenBlackByGin(c)
	if ok {
		slog.Warn("用户访问被拦截，token已在黑名单", "path", c.Request.URL.Path, "method", c.Request.Method, "client_ip", c.ClientIP(), "user_id", claims.UserID, "username", claims.Username, "black_type", blk.String())
		res.FailWithMsg(blk.Msg(), c)
		c.Abort()
		return
	}
	c.Set("claims", claims)
}
func AdminAuthMiddleware(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		slog.Warn("管理员鉴权失败", "path", c.Request.URL.Path, "method", c.Request.Method, "client_ip", c.ClientIP(), "err", err)
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	blk, ok := redis_jwt.HasTokenBlackByGin(c)
	if ok {
		slog.Warn("管理员访问被拦截，token已在黑名单", "path", c.Request.URL.Path, "method", c.Request.Method, "client_ip", c.ClientIP(), "user_id", claims.UserID, "username", claims.Username, "black_type", blk.String())
		res.FailWithMsg(blk.Msg(), c)
		c.Abort()
		return
	}
	if claims.Role != enum.AdminRole {
		slog.Warn("管理员接口权限不足", "path", c.Request.URL.Path, "method", c.Request.Method, "client_ip", c.ClientIP(), "user_id", claims.UserID, "username", claims.Username, "role", claims.Role)
		res.FailWithMsg("非管理员角色", c)
		c.Abort()
		return
	}
	c.Set("claims", claims)
}
