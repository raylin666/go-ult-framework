// Package middleware 提供基于 HTTP 框架的中间件管理系统。
package middleware

import (
	"github.com/gin-gonic/gin"
	utilsMiddleware "github.com/raylin666/go-utils/v2/middleware"
)

// Priority 中间件优先级类型。
type Priority = utilsMiddleware.Priority

// HandlerFunc 中间件处理函数类型。
type HandlerFunc func(ctx *gin.Context)

// Middleware HTTP 框架中间件接口。
// 扩展通用 Middleware 接口，提供框架特定的处理函数。
type Middleware interface {
	utilsMiddleware.Middleware
	// Handler 返回中间件处理函数。
	// 这是一个类型安全的接口，确保返回 HandlerFunc。
	Handler() utilsMiddleware.Handler
}

// middlewareFunc 函数式中间件实现。
type middlewareFunc struct {
	utilsMiddleware.Middleware
	handler HandlerFunc
}

// Handler 返回中间件处理函数。
func (m *middlewareFunc) Handler() utilsMiddleware.Handler {
	return m.handler
}

// NewMiddlewareFunc 创建函数式中间件。
//
// 参数:
//   - name: 中间件名称
//   - priority: 中间件优先级
//   - handler: 中间件处理函数
//
// 返回:
//   - Middleware: 中间件实例
func NewMiddlewareFunc(name string, priority utilsMiddleware.Priority, handler HandlerFunc) Middleware {
	base := utilsMiddleware.NewMiddlewareFunc(name, priority, handler)
	return &middlewareFunc{
		Middleware: base,
		handler:    handler,
	}
}
