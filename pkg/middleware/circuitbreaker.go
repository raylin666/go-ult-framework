// Package middleware 提供常用中间件实现。
// 包括熔断器（Circuit Breaker）和限流器（Rate Limiter）等。
package middleware

import (
	"errors"
	"sync"
	"time"
)

// 熔断器错误定义。
var (
	ErrCircuitBreakerOpen = errors.New("circuit breaker is open") // 熔断器打开错误
)

// CircuitState 熔断器状态类型。
type CircuitState int

// 熔断器状态常量。
const (
	StateClosed   CircuitState = iota // 关闭状态（正常）
	StateOpen                         // 打开状态（熔断）
	StateHalfOpen                     // 半开状态（恢复尝试）
)

// CircuitBreaker 熔断器实现。
// 用于防止级联故障，当失败次数达到阈值时自动熔断。
type CircuitBreaker struct {
	mu               sync.RWMutex  // 读写锁
	state            CircuitState  // 当前状态
	failureCount     int           // 失败计数
	successCount     int           // 成功计数
	failureThreshold int           // 失败阈值
	successThreshold int           // 成功阈值（半开状态下恢复所需）
	timeout          time.Duration // 熔断超时时间
	lastFailureTime  time.Time     // 最后失败时间
}

// CircuitBreakerOption 熔断器配置选项。
type CircuitBreakerOption struct {
	FailureThreshold int           // 失败阈值
	SuccessThreshold int           // 成功阈值
	Timeout          time.Duration // 熔断超时时间
}

// NewCircuitBreaker 创建新的熔断器实例。
// 默认失败阈值 5，成功阈值 2，超时时间 30 秒。
//
// 参数:
//   - opts: 熔断器配置选项
//
// 返回:
//   - *CircuitBreaker: 熔断器实例
func NewCircuitBreaker(opts CircuitBreakerOption) *CircuitBreaker {
	if opts.FailureThreshold <= 0 {
		opts.FailureThreshold = 5
	}
	if opts.SuccessThreshold <= 0 {
		opts.SuccessThreshold = 2
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 30 * time.Second
	}

	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: opts.FailureThreshold,
		successThreshold: opts.SuccessThreshold,
		timeout:          opts.Timeout,
	}
}

// State 获取当前熔断器状态。
//
// 返回:
//   - CircuitState: 当前状态
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Allow 检查是否允许请求通过。
// 在关闭状态下始终允许，在打开状态下超时后允许（进入半开）。
//
// 返回:
//   - bool: true 表示允许，false 表示拒绝
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailureTime) > cb.timeout {
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

// RecordSuccess 记录成功请求。
// 在关闭状态下重置失败计数，在半开状态下增加成功计数并可能恢复到关闭状态。
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		cb.failureCount = 0
	case StateHalfOpen:
		cb.successCount++
		if cb.successCount >= cb.successThreshold {
			cb.state = StateClosed
			cb.failureCount = 0
			cb.successCount = 0
		}
	case StateOpen:
		cb.state = StateHalfOpen
		cb.successCount = 1
	}
}

// RecordFailure 记录失败请求。
// 在关闭状态下增加失败计数并可能触发熔断，在半开状态下立即回到打开状态。
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.lastFailureTime = time.Now()

	switch cb.state {
	case StateClosed:
		cb.failureCount++
		if cb.failureCount >= cb.failureThreshold {
			cb.state = StateOpen
			cb.successCount = 0
		}
	case StateHalfOpen:
		cb.state = StateOpen
		cb.successCount = 0
	case StateOpen:
	}
}

// Execute 在熔断器保护下执行函数。
// 如果熔断器打开则返回错误，否则执行函数并记录结果。
//
// 参数:
//   - fn: 要执行的函数
//
// 返回:
//   - error: 执行错误或熔断器打开错误
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.Allow() {
		return ErrCircuitBreakerOpen
	}

	err := fn()
	if err != nil {
		cb.RecordFailure()
		return err
	}

	cb.RecordSuccess()
	return nil
}

// CircuitBreakerManager 熔断器管理器。
// 管理多个命名熔断器实例，支持按名称获取和执行。
type CircuitBreakerManager struct {
	mu       sync.RWMutex               // 读写锁
	breakers map[string]*CircuitBreaker // 熔断器映射
	opts     CircuitBreakerOption       // 默认配置选项
}

// NewCircuitBreakerManager 创建新的熔断器管理器。
//
// 参数:
//   - opts: 熔断器配置选项
//
// 返回:
//   - *CircuitBreakerManager: 熔断器管理器实例
func NewCircuitBreakerManager(opts CircuitBreakerOption) *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
		opts:     opts,
	}
}

// Get 根据名称获取熔断器实例。
// 如果不存在则创建新的熔断器。
//
// 参数:
//   - name: 熔断器名称
//
// 返回:
//   - *CircuitBreaker: 熔断器实例
func (m *CircuitBreakerManager) Get(name string) *CircuitBreaker {
	m.mu.RLock()
	if cb, exists := m.breakers[name]; exists {
		m.mu.RUnlock()
		return cb
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	cb := NewCircuitBreaker(m.opts)
	m.breakers[name] = cb
	return cb
}

// Execute 在指定名称的熔断器保护下执行函数。
//
// 参数:
//   - name: 熔断器名称
//   - fn: 要执行的函数
//
// 返回:
//   - error: 执行错误或熔断器打开错误
func (m *CircuitBreakerManager) Execute(name string, fn func() error) error {
	return m.Get(name).Execute(fn)
}
