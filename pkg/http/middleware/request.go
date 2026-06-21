// Package middleware 提供基于 HTTP 框架的中间件管理系统。
package middleware

import (
	"net/http"
	"time"

	"ult/pkg/types"
	pkgtypes "ult/pkg/types"

	"github.com/gin-gonic/gin"
	utilsMiddleware "github.com/raylin666/go-utils/v2/middleware"
	"github.com/raylin666/go-utils/v2/validator"
)

// RequestConfig Request 中间件配置。
type RequestConfig struct {
	// Enabled 是否启用请求处理中间件
	Enabled bool

	// Validator 数据验证器
	Validator validator.Validator

	// ContextInitializer Context 初始化函数
	// 该函数负责创建和初始化 Context，包括读取请求体、设置验证器等
	// 返回的 Context 应该已经完成初始化（调用过 init() 方法）
	ContextInitializer ContextInitializer

	// Response 响应处理函数
	// 该函数负责处理响应，包括错误响应、成功响应、日志记录等
	Response Response
}

// ContextInitializer Context 初始化函数类型。
// 该函数负责创建和初始化 Context，包括读取请求体、设置验证器等。
// 返回的 Context 应该已经完成初始化（调用过 init() 方法）。
type ContextInitializer func(ctx *gin.Context) (interface{}, error)

// ResponseHandler 响应处理函数类型。
// 该函数负责处理响应，包括错误响应、成功响应、日志记录等。
type Response func(reqTime time.Time, ctx *gin.Context)

// Request 请求处理中间件。
// 负责 Context 初始化、验证器设置和响应处理。
// 该中间件应该在所有其他中间件之前执行，确保 Context 在中间件中可用。
type Request struct {
	config *RequestConfig
}

// NewRequest 创建 Request 中间件实例。
//
// 参数:
//   - config: Request 配置
//
// 返回:
//   - *Request: Request 中间件实例
func NewRequest(config *RequestConfig) *Request {
	return &Request{
		config: config,
	}
}

// Name 返回中间件名称。
func (r *Request) Name() string {
	return types.RequestMiddlewareName
}

// Priority 返回中间件优先级。
// Request 必须在所有业务中间件之前执行，设置为高优先级。
// 但比 Recovery 稍低（Recovery 是 PriorityHighest=0，这里使用 PriorityHigh=1），
// 确保 Recovery 能捕获 Request 中的 panic。
func (r *Request) Priority() utilsMiddleware.Priority {
	return utilsMiddleware.PriorityHigh
}

// Enabled 返回是否启用。
func (r *Request) Enabled() bool {
	return r.config.Enabled
}

// Handler 返回中间件处理函数（实现 utilsMiddleware.Middleware 接口）。
func (r *Request) Handler() utilsMiddleware.Handler {
	return r.handler()
}

// handler 返回中间件处理函数。
// 初始化 Context、设置验证器、处理响应。
func (r *Request) handler() HandlerFunc {
	return func(ctx *gin.Context) {
		// 拦截 404 请求路由
		if ctx.Writer.Status() == http.StatusNotFound {
			return
		}

		// 请求时间
		ts := time.Now()

		// 初始化 Context（通过配置的初始化函数）
		if r.config.ContextInitializer != nil {
			appCtx, err := r.config.ContextInitializer(ctx)
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			// 存储 Context 到 gin.Context（appCtx 是 interface{} 类型）
			ctx.Set(pkgtypes.CoreContextNameKey, appCtx)
		}

		// 响应处理（通过配置的响应处理函数）
		if r.config.Response != nil {
			defer func() {
				r.config.Response(ts, ctx)
			}()
		}

		ctx.Next()
	}
}

// DefaultRequestConfig 返回默认 Request 配置。
//
// 参数:
//   - validator: 数据验证器
//   - contextInitializer: Context 初始化函数
//   - response: 响应处理函数
//
// 返回:
//   - *RequestConfig: 默认 Request 配置
func DefaultRequestConfig(
	validator validator.Validator,
	contextInitializer ContextInitializer,
	response Response,
) *RequestConfig {
	return &RequestConfig{
		Enabled:            true,
		Validator:          validator,
		ContextInitializer: contextInitializer,
		Response:           response,
	}
}

// NewDefaultRequest 创建默认 Request 中间件。
// 提供便捷的创建方式，使用默认配置。
//
// 参数:
//   - validator: 数据验证器
//   - contextInitializer: Context 初始化函数
//   - response: 响应处理函数
//
// 返回:
//   - *Request: Request 中间件实例
func NewDefaultRequest(
	validator validator.Validator,
	contextInitializer ContextInitializer,
	response Response,
) *Request {
	return NewRequest(DefaultRequestConfig(validator, contextInitializer, response))
}
