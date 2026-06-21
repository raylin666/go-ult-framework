// Package middleware 提供基于 HTTP 框架的中间件管理系统。
package middleware

import (
	"fmt"

	utilsMiddleware "github.com/raylin666/go-utils/v2/middleware"
)

// Manager 中间件管理器。
// 包装通用 Manager，提供框架特定的方法。
type Manager struct {
	*utilsMiddleware.Manager
	middlewares []Middleware // 存储中间件实例，用于依赖验证
}

// NewManager 创建中间件管理器。
//
// 返回:
//   - *Manager: 中间件管理器实例
func NewManager() *Manager {
	return &Manager{
		Manager:     utilsMiddleware.NewManager(),
		middlewares: []Middleware{},
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
	m.middlewares = append(m.middlewares, middleware) // 存储实例用于依赖验证
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
// 按优先级排序，验证依赖关系，返回 HandlerFunc 列表。
//
// 返回:
//   - []HandlerFunc: 中间件处理函数列表
func (m *Manager) Build() []HandlerFunc {
	// 验证中间件依赖关系（使用存储的 Middleware 实例）
	m.validateDependencies()

	handlers := m.Manager.Build()
	httpHandlers := make([]HandlerFunc, 0, len(handlers))

	for _, h := range handlers {
		// 类型断言，确保是 HandlerFunc
		if httpHandler, ok := h.(HandlerFunc); ok {
			httpHandlers = append(httpHandlers, httpHandler)
		}
	}

	return httpHandlers
}

// validateDependencies 验证中间件依赖关系。
// 检查每个中间件的依赖是否都已注册。
func (m *Manager) validateDependencies() {
	// 构建已注册中间件的名称集合
	registered := make(map[string]bool)
	for _, middleware := range m.middlewares {
		registered[middleware.Name()] = true
	}

	// 验证每个中间件的依赖
	for _, middleware := range m.middlewares {
		for _, dep := range middleware.Dependencies() {
			if !registered[dep] {
				panic(fmt.Sprintf(
					"中间件依赖验证失败: 中间件 '%s' 依赖 '%s'，但 '%s' 未注册",
					middleware.Name(),
					dep,
					dep,
				))
			}
		}
	}
}
