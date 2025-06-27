package middleware

import (
	"bytes"
	"io"
	"time"

	"MetaFarmBackend/component/logger"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 复制请求体和响应体
		var requestBody bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &requestBody)
		c.Request.Body = io.NopCloser(tee)

		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = writer

		// 处理请求
		c.Next()

		// 记录请求信息
		latency := time.Since(startTime)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()
		token := c.GetHeader("Authorization")
		contentType := c.GetHeader("Content-Type")

		logger.Infof(
			"Request | %3d | %13v | %-7s %s | %15s | %s | Token: %s | Content-Type: %s | Request: %s | Response: %s",
			status,
			latency,
			method,
			path,
			ip,
			userAgent,
			token,
			contentType,
			requestBody.String(),
			writer.body.String(),
		)
	}
}
