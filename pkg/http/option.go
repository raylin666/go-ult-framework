// Package http 提供 HTTP 服务器实现，基于 Gin 框架封装。
package http

import (
	"time"

	pkgmiddleware "ult/pkg/http/middleware"
)

// Option HTTP 服务器选项函数类型。
type Option func(opt *option)

// option HTTP 服务器内部选项配置。
type option struct {
	openBrowser string                     // 启动时自动打开的浏览器 URL
	timeout     time.Duration              // 优雅关闭超时时间
	middlewares []pkgmiddleware.Middleware // 自定义中间件列表
}

// EnableOpenBrowser 启动时自动打开浏览器选项。
//
// 参数:
//   - uri: 要打开的 URL
//
// 返回:
//   - Option: 选项函数
func EnableOpenBrowser(uri string) Option {
	return func(opt *option) {
		opt.openBrowser = uri
	}
}

// WithTimeout 设置优雅关闭超时时间选项。
//
// 参数:
//   - ts: 超时时间
//
// 返回:
//   - Option: 选项函数
func WithTimeout(ts time.Duration) Option {
	return func(opt *option) {
		opt.timeout = ts
	}
}

// WithMiddleware 添加自定义中间件选项。
// 支持新的中间件管理系统，提供优先级控制和灵活配置。
//
// 参数:
//   - m: 中间件列表
//
// 返回:
//   - Option: 选项函数
func WithMiddleware(m ...pkgmiddleware.Middleware) Option {
	return func(opt *option) {
		opt.middlewares = append(opt.middlewares, m...)
	}
}
