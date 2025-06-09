package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

// RateLimiter 限流中间件
func RateLimiter(limit int) gin.HandlerFunc {
	limiter := ratelimit.New(limit, ratelimit.Per(time.Second))
	return func(c *gin.Context) {
		limiter.Take()
		c.Next()
	}
}

// IPRateLimiter IP限流中间件
func IPRateLimiter(limit int) gin.HandlerFunc {
	limiters := make(map[string]ratelimit.Limiter)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter, exists := limiters[ip]
		if !exists {
			limiter = ratelimit.New(limit, ratelimit.Per(time.Second))
			limiters[ip] = limiter
		}
		limiter.Take()
		c.Next()
	}
}
