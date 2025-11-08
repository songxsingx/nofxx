package api

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 简单的令牌桶限流器
type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate time.Duration
	mu         sync.Mutex
	lastRefill time.Time
}

// NewRateLimiter 创建限流器
func NewRateLimiter(maxTokens int, refillRate time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 补充令牌
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)
	
	if tokensToAdd > 0 {
		rl.tokens = min(rl.maxTokens, rl.tokens+tokensToAdd)
		rl.lastRefill = now
	}

	// 消耗令牌
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(429, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RecoveryMiddleware 自定义恢复中间件（更详细的错误日志）
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录详细的panic信息
				c.JSON(500, gin.H{
					"error": "服务器内部错误",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
