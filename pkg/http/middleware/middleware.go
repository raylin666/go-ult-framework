// Package middleware 提供中间件管理系统。
// 该包定义了中间件接口、优先级管理、中间件链构建等功能。
package middleware

import (
	"github.com/gin-gonic/gin"
)

// HandlerFunc 中间件处理函数类型。
// 定义为 gin.HandlerFunc 以避免导入循环。
type HandlerFunc func(ctx *gin.Context)

// Priority 中间件优先级类型。
// 数值越小优先级越高，中间件按优先级顺序执行。
type Priority int

const (
	// PriorityHighest 最高优先级（数值最小）。
	// 用于必须在最前执行的中间件，如异常恢复、链路追踪。
	PriorityHighest Priority = iota

	// PriorityHigh 高优先级。
	// 用于需要在早期执行的中间件，如 CORS、安全检查。
	PriorityHigh

	// PriorityNormal 正常优先级。
	// 用于常规中间件，如日志、验证、请求处理。
	PriorityNormal

	// PriorityLow 低优先级。
	//用于业务相关的中间件，如权限检查、限流。
	PriorityLow
)

// Middleware 中间件接口。
// 所有中间件必须实现此接口，提供名称、优先级和处理函数。
type Middleware interface {
	// Name 返回中间件名称。
	// 用于日志记录、调试和中间件识别。
	Name() string

	// Priority 返回中间件优先级。
	// 决定中间件在链中的执行顺序。
	Priority() Priority

	// Handler 返回中间件处理函数。
	// 该函数将在 HTTP 请求处理过程中执行。
	Handler() HandlerFunc
}

// Config 中间件配置接口。
// 用于可配置的中间件，提供启用状态检查。
type Config interface {
	// Enabled 返回中间件是否启用。
	// 如果返回 false，中间件管理器将跳过该中间件。
	Enabled() bool
}

// middlewareFunc 函数式中间件实现。
// 提供简化的中间件创建方式，无需定义完整结构体。
type middlewareFunc struct {
	name     string
	priority Priority
	handler  HandlerFunc
}

// Name 返回中间件名称。
func (m *middlewareFunc) Name() string {
	return m.name
}

// Priority 返回中间件优先级。
func (m *middlewareFunc) Priority() Priority {
	return m.priority
}

// Handler 返回中间件处理函数。
func (m *middlewareFunc) Handler() HandlerFunc {
	return m.handler
}

// NewMiddlewareFunc 创建函数式中间件。
// 提供简化的中间件创建方式。
//
// 参数:
//   - name: 中间件名称
//   - priority: 中间件优先级
//   - handler: 中间件处理函数
//
// 返回:
//   - Middleware: 中间件实例
func NewMiddlewareFunc(name string, priority Priority, handler HandlerFunc) Middleware {
	return &middlewareFunc{
		name:     name,
		priority: priority,
		handler:  handler,
	}
}