// Package middleware 提供基于 HTTP 框架的中间件管理系统。
package middleware

import (
	"github.com/gin-gonic/gin"
	utilsMiddleware "github.com/raylin666/go-utils/v2/middleware"
)

// Manager 中间件管理器。
// 包装通用 Manager，提供框架特定的方法。
type Manager struct {
	*utilsMiddleware.Manager
}

// NewManager 创建中间件管理器。
//
// 返回:
//   - *Manager: 中间件管理器实例
func NewManager() *Manager {
	return &Manager{
		Manager: utilsMiddleware.NewManager(),
	}
}

// Use 添加中间件。
//
// 参数:
//   - middleware: 中间件实例
//
// 返回:
//   - *Manager: 中间件管理器实例（支持链式调用）
func (m *Manager) Use(middleware Middleware) *Manager {
	m.Manager.Use(middleware)
	return m
}

// UseFunc 使用函数方式添加中间件。
//
// 参数:
//   - name: 中间件名称
//   - priority: 中间件优先级
//   - handler: 中间件处理函数
//
// 返回:
//   - *Manager: 中间件管理器实例（支持链式调用）
func (m *Manager) UseFunc(name string, priority utilsMiddleware.Priority, handler HandlerFunc) *Manager {
	return m.Use(NewMiddlewareFunc(name, priority, handler))
}

// Build 构建中间件链。
// 按优先级排序，返回 HandlerFunc 列表。
//
// 返回:
//   - []HandlerFunc: 中间件处理函数列表
func (m *Manager) Build() []HandlerFunc {
	handlers := m.Manager.Build()
	httpHandlers := make([]HandlerFunc, len(handlers))

	for i, h := range handlers {
		// 类型断言，确保是 HandlerFunc
		if httpHandler, ok := h.(HandlerFunc); ok {
			httpHandlers[i] = httpHandler
		} else {
			// 如果不是 HandlerFunc，创建一个空处理函数
			httpHandlers[i] = func(ctx *gin.Context) {
				ctx.Next()
			}
		}
	}

	return httpHandlers
}
