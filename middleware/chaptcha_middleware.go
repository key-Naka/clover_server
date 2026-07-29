package middleware

import (
	"bytes"
	"clover_server/common/res"
	"clover_server/global"
	"io"

	"github.com/gin-gonic/gin"
)

type CaptchaMiddlewareRequest struct {
	Code string `json:"code"`
	ID   string `json:"id"`
}

func CaptchaMiddleware(c *gin.Context) {
	if !global.Config.Site.Captcha.Enable {
		c.Next()
		return
	}
	body, err := c.GetRawData()
	if err != nil {
		res.FailWithMsg("获取请求体错误", c)
		c.Abort()
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	var req CaptchaMiddlewareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("图形验证码解析错误", c)
		c.Abort()
		return
	}
	if !global.CaptchaStore.Verify(req.ID, req.Code, true) {
		res.FailWithMsg("验证码错误", c)
		c.Abort()
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

}
