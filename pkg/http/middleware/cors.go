// Package middleware 提供基于 HTTP 框架的中间件管理系统。
package middleware

import (
	nethttp "net/http"
	"ult/pkg/types"

	"github.com/gin-gonic/gin"
	utilsMiddleware "github.com/raylin666/go-utils/v2/middleware"
	cors "github.com/rs/cors/wrapper/gin"
)

// CORSConfig CORS 中间件配置。
type CORSConfig struct {
	// Enabled 是否启用 CORS 中间件
	Enabled bool

	// AllowedOrigins 允许的域名列表
	// 例如: ["https://example.com", "http://localhost:3000"] 或 ["*"]
	AllowedOrigins []string

	// AllowedMethods 允许的 HTTP 方法
	// 默认包含常用的 HTTP 方法
	AllowedMethods []string

	// AllowedHeaders 允许的请求头
	AllowedHeaders []string

	// AllowCredentials 是否允许携带凭证（如 Cookie）
	AllowCredentials bool

	// OptionsPassthrough 是否让 OPTIONS 请求继续传递
	OptionsPassthrough bool
}

// CORS CORS 中间件。
// 处理跨域资源共享（Cross-Origin Resource Sharing）。
type CORS struct {
	config       *CORSConfig
	corsHandler  gin.HandlerFunc // 缓存的 CORS 处理器，避免每次请求重新创建
}

// NewCORS 创建 CORS 中间件实例。
//
// 参数:
//   - config: CORS 配置
//
// 返回:
//   - *CORS: CORS 中间件实例
func NewCORS(config *CORSConfig) *CORS {
	// 设置默认值
	if config.AllowedMethods == nil {
		config.AllowedMethods = []string{
			nethttp.MethodHead,
			nethttp.MethodGet,
			nethttp.MethodPost,
			nethttp.MethodPut,
			nethttp.MethodPatch,
			nethttp.MethodDelete,
		}
	}

	// 在初始化时创建 CORS 处理器，避免每次请求重复创建
	var corsHandler gin.HandlerFunc
	if len(config.AllowedOrigins) > 0 {
		corsHandler = cors.New(cors.Options{
			AllowedOrigins:     config.AllowedOrigins,
			AllowedMethods:     config.AllowedMethods,
			AllowedHeaders:     config.AllowedHeaders,
			AllowCredentials:   config.AllowCredentials,
			OptionsPassthrough: config.OptionsPassthrough,
		})
	}

	return &CORS{
		config:       config,
		corsHandler:  corsHandler,
	}
}

// Name 返回中间件名称。
func (c *CORS) Name() string {
	return types.CorsMiddlewareName
}

// Priority 返回中间件优先级。
// CORS 中间件需要在早期执行，设置为高优先级。
func (c *CORS) Priority() utilsMiddleware.Priority {
	return utilsMiddleware.PriorityHigh
}

// Enabled 返回是否启用。
func (c *CORS) Enabled() bool {
	return c.config.Enabled
}

// Dependencies 返回中间件依赖列表。
// CORS 中间件无依赖。
func (c *CORS) Dependencies() []string {
	return []string{}
}

// Handler 返回中间件处理函数（实现 utilsMiddleware.Middleware 接口）。
func (c *CORS) Handler() utilsMiddleware.Handler {
	return c.handler()
}

// handler 返回中间件处理函数。
// 使用缓存的 CORS 处理器，避免每次请求重复创建。
func (c *CORS) handler() HandlerFunc {
	return func(ctx *gin.Context) {
		// 如果没有配置允许的域名或处理器未创建，跳过 CORS 处理
		if c.corsHandler == nil {
			ctx.Next()
			return
		}

		// 使用已缓存的 CORS 处理器
		c.corsHandler(ctx)
	}
}

// DefaultCORSConfig 返回默认 CORS 配置。
// 采用安全配置，默认不允许任何域名，强制用户显式配置。
//
// 返回:
//   - *CORSConfig: 默认 CORS 配置
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		Enabled:        true,
		AllowedOrigins: []string{}, // 默认不允许任何域名，强制用户配置
		AllowedMethods: []string{
			nethttp.MethodHead,
			nethttp.MethodGet,
			nethttp.MethodPost,
			nethttp.MethodPut,
			nethttp.MethodPatch,
			nethttp.MethodDelete,
		},
		AllowCredentials:   false,
		OptionsPassthrough: false,
	}
}
