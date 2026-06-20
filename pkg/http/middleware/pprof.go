// Package middleware 提供基于 HTTP 框架的中间件管理系统。
package middleware

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	utilsMiddleware "github.com/raylin666/go-utils/v2/middleware"
	"github.com/raylin666/go-utils/v2/server/system"
)

// PProfConfig PProf 中间件配置。
type PProfConfig struct {
	// Enabled 是否启用 PProf 性能分析
	Enabled bool

	// Environment 环境标识（用于判断是否在生产环境）
	Environment string
}

// PProf 性能分析中间件。
// 提供 Go 程序的性能分析功能，访问路径: /debug/pprof
type PProf struct {
	config *PProfConfig
}

// NewPProf 创建 PProf 中间件实例。
//
// 参数:
//   - config: PProf 配置
//
// 返回:
//   - *PProf: PProf 中间件实例
func NewPProf(config *PProfConfig) *PProf {
	return &PProf{
		config: config,
	}
}

// Name 返回中间件名称。
func (p *PProf) Name() string {
	return "pprof"
}

// Priority 返回中间件优先级。
// PProf 中间件优先级较低，在业务逻辑之后执行。
func (p *PProf) Priority() utilsMiddleware.Priority {
	return utilsMiddleware.PriorityLow
}

// Enabled 返回是否启用。
// 只有在非生产环境且配置启用时才真正启用。
func (p *PProf) Enabled() bool {
	if !p.config.Enabled {
		return false
	}
	// 在生产环境不启用 PProf
	return !system.NewEnvironment(p.config.Environment).IsProd()
}

// Handler 返回中间件处理函数（实现 utilsMiddleware.Middleware 接口）。
func (p *PProf) Handler() utilsMiddleware.Handler {
	return p.handler()
}

// handler 返回中间件处理函数。
// 注册 pprof 路由到 Gin 引擎。
func (p *PProf) handler() HandlerFunc {
	return func(ctx *gin.Context) {
		// PProf 通过 pprof.Register 注册路由，不需要在中间件中处理
		// 这里只是一个占位符，实际的注册逻辑在 server.go 中处理
		ctx.Next()
	}
}

// RegisterRoutes 注册 PProf 路由到 Gin 引擎。
// 这是一个特殊方法，因为 PProf 需要注册路由而不是作为中间件使用。
//
// 参数:
//   - engine: Gin 引擎实例
func (p *PProf) RegisterRoutes(engine *gin.Engine) {
	if p.Enabled() {
		pprof.Register(engine)
	}
}

// DefaultPProfConfig 返回默认 PProf 配置。
// 默认在非生产环境启用。
//
// 参数:
//   - environment: 环境标识
//
// 返回:
//   - *PProfConfig: 默认 PProf 配置
func DefaultPProfConfig(environment string) *PProfConfig {
	return &PProfConfig{
		Enabled:     true,
		Environment: environment,
	}
}