// Package http 提供 HTTP 服务器实现，基于 Gin 框架封装。
package http

import (
	"time"
	"ult/pkg/proposal"

	pkgmiddleware "ult/pkg/http/middleware"
)

// Option HTTP 服务器选项函数类型。
type Option func(opt *option)

// option HTTP 服务器内部选项配置。
type option struct {
	cors struct {
		domains []string // CORS 允许的域名列表
	}
	pprof       bool                       // 是否启用 pprof 性能分析
	rate        bool                       // 是否启用限流
	openBrowser string                     // 启动时自动打开的浏览器 URL
	alertNotify proposal.NotifyHandler     // 告警通知处理函数
	timeout     time.Duration              // 优雅关闭超时时间
	middlewares []pkgmiddleware.Middleware // 自定义中间件列表
}

// EnableCors 启用 CORS 跨域支持选项。
//
// 参数:
//   - domains: 允许跨域的域名列表
//
// 返回:
//   - Option: 选项函数
func EnableCors(domains []string) Option {
	return func(opt *option) {
		opt.cors.domains = domains
	}
}

// EnablePProf 启用 pprof 性能分析选项。
//
// 返回:
//   - Option: 选项函数
func EnablePProf() Option {
	return func(opt *option) {
		opt.pprof = true
	}
}

// EnableRate 启用限流选项。
//
// 返回:
//   - Option: 选项函数
func EnableRate() Option {
	return func(opt *option) {
		opt.rate = true
	}
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

// EnableAlertNotify 启用告警通知选项。
//
// 参数:
//   - handler: 告警通知处理函数
//
// 返回:
//   - Option: 选项函数
func EnableAlertNotify(handler proposal.NotifyHandler) Option {
	return func(opt *option) {
		opt.alertNotify = handler
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
