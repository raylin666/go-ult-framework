// Package middleware 提供常用中间件实现。
package middleware

import (
	stdhttp "net/http"
	"sync"
	"time"
	"ult/errcode"
	"ult/pkg/http"

	"github.com/gin-gonic/gin"
)

// 限流器默认配置常量。
const (
	RateLimitDefaultRequests = 100         // 默认最大请求数
	RateLimitDefaultWindow   = time.Second // 默认时间窗口
)

// RateLimiter 限流器实现。
// 基于 IP 地址进行限流，使用滑动窗口算法。
type RateLimiter struct {
	mu          sync.RWMutex           // 读写锁
	requests    map[string]*clientInfo // 客户端请求信息映射
	maxRequests int                    // 最大请求数
	window      time.Duration          // 时间窗口
}

// clientInfo 客户端请求信息。
type clientInfo struct {
	count     int       // 请求计数
	startTime time.Time // 窗口开始时间
}

// NewRateLimiter 创建新的限流器实例。
// 默认最大请求数 100，时间窗口 1 秒。
// 启动后台清理协程定期清理过期客户端信息。
//
// 参数:
//   - maxRequests: 最大请求数
//   - window: 时间窗口
//
// 返回:
//   - *RateLimiter: 限流器实例
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	if maxRequests <= 0 {
		maxRequests = RateLimitDefaultRequests
	}
	if window <= 0 {
		window = RateLimitDefaultWindow
	}

	limiter := &RateLimiter{
		requests:    make(map[string]*clientInfo),
		maxRequests: maxRequests,
		window:      window,
	}

	go limiter.cleanupRoutine()

	return limiter
}

// cleanupRoutine 后台清理协程。
// 每分钟清理超过两倍时间窗口未活跃的客户端信息。
func (rl *RateLimiter) cleanupRoutine() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, info := range rl.requests {
			if now.Sub(info.startTime) > rl.window*2 {
				delete(rl.requests, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow 检查指定 IP 是否允许请求。
// 使用滑动窗口算法，在时间窗口内超过最大请求数则拒绝。
//
// 参数:
//   - ip: 客户端 IP 地址
//
// 返回:
//   - bool: true 表示允许，false 表示拒绝
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	info, exists := rl.requests[ip]
	if !exists {
		rl.requests[ip] = &clientInfo{
			count:     1,
			startTime: now,
		}
		return true
	}

	if now.Sub(info.startTime) > rl.window {
		info.count = 1
		info.startTime = now
		return true
	}

	if info.count >= rl.maxRequests {
		return false
	}

	info.count++
	return true
}

// RateLimit Gin 限流中间件。
// 检查客户端 IP 是否超过限流阈值，超过则返回 429 错误。
//
// 参数:
//   - limiter: 限流器实例
//
// 返回:
//   - gin.HandlerFunc: Gin 中间件函数
func RateLimit(limiter *RateLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		if !limiter.Allow(ip) {
			appCtx, ok := ctx.Value(http.CoreContextNameKey).(http.Context)
			if ok && appCtx != nil {
				appCtx.WithAbortError(errcode.New(errcode.RequestError).WithDesc("rate limit exceeded"))
			} else {
				ctx.AbortWithStatusJSON(stdhttp.StatusTooManyRequests, gin.H{
					"code":    errcode.RequestError,
					"message": "rate limit exceeded",
				})
			}
			return
		}

		ctx.Next()
	}
}
