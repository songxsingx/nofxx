package exchange

import (
	"fmt"
	"log"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries   int           // 最大重试次数
	InitialDelay time.Duration // 初始延迟
	MaxDelay     time.Duration // 最大延迟
	Multiplier   float64       // 延迟倍增系数
}

// DefaultRetryConfig 默认重试配置
var DefaultRetryConfig = RetryConfig{
	MaxRetries:   3,
	InitialDelay: 1 * time.Second,
	MaxDelay:     10 * time.Second,
	Multiplier:   2.0,
}

// WithRetry 通用重试包装器
func WithRetry[T any](operation func() (T, error), config RetryConfig, operationName string) (T, error) {
	var result T
	var err error
	delay := config.InitialDelay

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		result, err = operation()
		if err == nil {
			if attempt > 0 {
				log.Printf("✓ %s 重试成功 (尝试 %d/%d)", operationName, attempt+1, config.MaxRetries+1)
			}
			return result, nil
		}

		if attempt < config.MaxRetries {
			log.Printf("⚠️ %s 失败，%v 后重试 (%d/%d): %v", 
				operationName, delay, attempt+1, config.MaxRetries, err)
			time.Sleep(delay)
			
			// 指数退避
			delay = time.Duration(float64(delay) * config.Multiplier)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		}
	}

	return result, fmt.Errorf("%s 重试 %d 次后失败: %w", operationName, config.MaxRetries, err)
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	maxFailures  int           // 最大失败次数
	resetTimeout time.Duration // 重置超时
	failures     int           // 当前失败次数
	lastFailTime time.Time     // 上次失败时间
	state        string        // 状态: closed, open, half-open
}

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        "closed",
	}
}

// Call 通过熔断器调用函数
func (cb *CircuitBreaker) Call(operation func() error, operationName string) error {
	// 检查是否需要从open状态转为half-open
	if cb.state == "open" {
		if time.Since(cb.lastFailTime) > cb.resetTimeout {
			log.Printf("🔧 熔断器半开: %s", operationName)
			cb.state = "half-open"
			cb.failures = 0
		} else {
			return fmt.Errorf("熔断器开启，拒绝执行 %s", operationName)
		}
	}

	// 执行操作
	err := operation()
	
	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()
		
		if cb.failures >= cb.maxFailures {
			log.Printf("🚨 熔断器打开: %s (失败次数: %d)", operationName, cb.failures)
			cb.state = "open"
		}
		return err
	}

	// 成功则重置
	if cb.state == "half-open" {
		log.Printf("✓ 熔断器关闭: %s", operationName)
		cb.state = "closed"
	}
	cb.failures = 0
	return nil
}

// GetState 获取熔断器状态
func (cb *CircuitBreaker) GetState() string {
	return cb.state
}

// Reset 重置熔断器
func (cb *CircuitBreaker) Reset() {
	cb.failures = 0
	cb.state = "closed"
}
