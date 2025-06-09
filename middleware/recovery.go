// Package middleware 提供了EasyGo框架的常用中间件
package middleware

import (
	"github.com/xzl-go/easygo/core"
	"github.com/xzl-go/easygo/logger"
)

// Recovery 返回一个恢复中间件
func Recovery() core.HandlerFunc {
	return func(c *core.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered: %v", err)
				c.JSON(500, map[string]string{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}
