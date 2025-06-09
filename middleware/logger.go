// Package middleware 提供了EasyGo框架的常用中间件
package middleware

import (
	"time"

	"github.com/xzl-go/easygo/core"
	"github.com/xzl-go/easygo/logger"
)

// Logger 返回一个日志中间件
func Logger() core.HandlerFunc {
	return func(c *core.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.Request.RemoteAddr
		method := c.Request.Method
		statusCode := c.StatusCode

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("[%s] %s %s %d %v",
			clientIP,
			method,
			path,
			statusCode,
			latency,
		)
	}
}
