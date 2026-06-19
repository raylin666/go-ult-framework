// Package middleware 提供中间件管理系统。
package middleware

import (
	nethttp "net/http"

	cors "github.com/rs/cors/wrapper/gin"
	"github.com/gin-gonic/gin"
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
	config *CORSConfig
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

	return &CORS{
		config: config,
	}
}

// Name 返回中间件名称。
func (c *CORS) Name() string {
	return "cors"
}

// Priority 返回中间件优先级。
// CORS 中间件需要在早期执行，设置为高优先级。
func (c *CORS) Priority() Priority {
	return PriorityHigh
}

// Enabled 返回是否启用。
func (c *CORS) Enabled() bool {
	return c.config.Enabled
}

// Handler 返回中间件处理函数。
// 使用 rs/cors 库实现 CORS 处理。
func (c *CORS) Handler() HandlerFunc {
	return func(ctx *gin.Context) {
		// 如果没有配置允许的域名，跳过 CORS 处理
		if len(c.config.AllowedOrigins) == 0 {
			return
		}

		// 创建 CORS 处理器
		corsHandler := cors.New(cors.Options{
			AllowedOrigins:     c.config.AllowedOrigins,
			AllowedMethods:     c.config.AllowedMethods,
			AllowedHeaders:     c.config.AllowedHeaders,
			AllowCredentials:   c.config.AllowCredentials,
			OptionsPassthrough: c.config.OptionsPassthrough,
		})

		// 执行 CORS 处理
		corsHandler(ctx)
	}
}

// DefaultCORSConfig 返回默认 CORS 配置。
// 允许所有域名访问，适用于开发环境。
//
// 返回:
//   - *CORSConfig: 默认 CORS 配置
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		Enabled:           true,
		AllowedOrigins:    []string{"*"},
		AllowedMethods: []string{
			nethttp.MethodHead,
			nethttp.MethodGet,
			nethttp.MethodPost,
			nethttp.MethodPut,
			nethttp.MethodPatch,
			nethttp.MethodDelete,
		},
		AllowCredentials:   true,
		OptionsPassthrough: true,
	}
}