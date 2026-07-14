package middleware

import (
	"clover_server/common/res"
	"clover_server/models/enum"
	"clover_server/service/redis_service/redis_jwt"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	blk, ok := redis_jwt.HasTokenBlackByGin(c)
	if ok {
		res.FailWithMsg(blk.Msg(), c)
		c.Abort()
		return
	}
	c.Set("claims", claims)
	return
}
func AdminAuthMiddleware(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	blk, ok := redis_jwt.HasTokenBlackByGin(c)
	if ok {
		res.FailWithMsg(blk.Msg(), c)
		c.Abort()
		return
	}
	if claims.Role != enum.AdminRole {
		res.FailWithMsg("非管理员角色", c)
		c.Abort()
		return
	}
	c.Set("claims", claims)
}
