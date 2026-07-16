package middleware

import (
	"bytes"
	"io"
	"log/slog"

	"clover_server/service/log_service"

	"github.com/gin-gonic/gin"
)

type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *ResponseWriter) Write(p []byte) (n int, err error) {
	w.body.Write(p)
	return w.ResponseWriter.Write(p)
}
func LogMiddleware(c *gin.Context) {
	ActionLog := log_service.GetLog(c)

	var bodyBytes []byte
	if c.Request.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(c.Request.Body)
		if err != nil {
			slog.Warn("读取请求体失败", "path", c.Request.URL.Path, "method", c.Request.Method, "client_ip", c.ClientIP(), "err", err)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	ActionLog.SetRequestBody(bodyBytes)
	c.Set("log", ActionLog)
	resw := &ResponseWriter{
		ResponseWriter: c.Writer,
		body:           bytes.NewBufferString(""),
	}
	c.Writer = resw
	c.Next()
	ActionLog.SetResponseBody(resw.body.Bytes())
	ActionLog.SetResponseHeader(c.Writer.Header())
	ActionLog.MiddlewareSave()
}
