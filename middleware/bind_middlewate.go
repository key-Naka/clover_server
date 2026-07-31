package middleware

import (
	"clover_server/common/res"

	"github.com/gin-gonic/gin"
)

func BindJSONMiddleware[T any](c *gin.Context) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	c.Set("req", req)
	return
}
func BindQueryMiddleware[T any](c *gin.Context) {
	var req T
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	c.Set("req", req)
	return
}
func BindUriMiddleware[T any](c *gin.Context) {
	var req T
	if err := c.ShouldBindUri(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	c.Set("req", req)
	return
}
func GetReq[T any](c *gin.Context) T {
	return c.MustGet("req").(T)
}
