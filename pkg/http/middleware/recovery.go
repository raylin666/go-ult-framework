// Package middleware 提供基于 HTTP 框架的中间件管理系统。
package middleware

import (
	"fmt"
	nethttp "net/http"
	"runtime/debug"
	"time"

	"ult/config"
	"ult/errcode"
	"ult/pkg/logger"
	"ult/pkg/proposal"
	"ult/pkg/types"
	pkgtypes "ult/pkg/types"

	goerror "errors"

	"github.com/gin-gonic/gin"
	utilsMiddleware "github.com/raylin666/go-utils/v2/middleware"
	"go.uber.org/zap"
)

// RecoveryConfig Recovery 中间件配置。
type RecoveryConfig struct {
	// Enabled 是否启用异常恢复中间件
	Enabled bool

	// AlertNotify 告警通知处理函数
	// 当发生 panic 时，会调用此函数发送告警通知
	AlertNotify proposal.NotifyHandler

	// Config 应用配置（用于告警通知）
	Config *config.Config

	// PrintStack 是否打印堆栈信息
	PrintStack bool
}

// Recovery 异常恢复中间件。
// 捕获 panic 并进行恢复，记录错误日志，发送告警通知。
type Recovery struct {
	config *RecoveryConfig
	logger *logger.Logger
}

// NewRecovery 创建 Recovery 中间件实例。
//
// 参数:
//   - config: Recovery 配置
//   - logger: 日志记录器
//
// 返回:
//   - *Recovery: Recovery 中间件实例
func NewRecovery(config *RecoveryConfig, logger *logger.Logger) *Recovery {
	return &Recovery{
		config: config,
		logger: logger,
	}
}

// Name 返回中间件名称。
func (r *Recovery) Name() string {
	return types.RecoveryMiddlewareName
}

// Priority 返回中间件优先级。
// Recovery 中间件必须在最前执行，设置为最高优先级。
func (r *Recovery) Priority() utilsMiddleware.Priority {
	return utilsMiddleware.PriorityHighest
}

// Enabled 返回是否启用。
func (r *Recovery) Enabled() bool {
	return r.config.Enabled
}

// Dependencies 返回中间件依赖列表。
// Recovery 中间件无依赖，必须最先执行。
func (r *Recovery) Dependencies() []string {
	return []string{}
}

// Handler 返回中间件处理函数（实现 utilsMiddleware.Middleware 接口）。
func (r *Recovery) Handler() utilsMiddleware.Handler {
	return r.handler()
}

// handler 返回中间件处理函数。
// 捕获 panic 并进行恢复处理。
func (r *Recovery) handler() HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				r.handleRecovery(ctx, err)
			}
		}()
		ctx.Next()
	}
}

// handleRecovery 处理 panic 恢复逻辑。
// 记录错误日志、堆栈信息，发送告警通知，设置错误响应。
//
// 参数:
//   - ctx: Gin 上下文
//   - err: panic 错误信息
func (r *Recovery) handleRecovery(ctx *gin.Context, err interface{}) {
	// 获取堆栈信息
	var stack = string(debug.Stack())

	// 记录错误日志
	r.logger.UseApp(ctx).Error(
		"got panic",
		zap.String("panic", fmt.Sprintf("%+v", err)),
		zap.String("stack", stack),
	)

	// 设置错误信息到 Context（先设置错误，再终止请求）
	ctx.Set(pkgtypes.ContextAbortErrorNameKey, errcode.New(errcode.ServerError).WithStackError(goerror.New("got panic")))

	// 终止请求并设置响应状态码
	ctx.AbortWithStatus(nethttp.StatusInternalServerError)

	// 发送告警通知
	if r.config.AlertNotify != nil && r.config.Config != nil {
		r.config.AlertNotify(&proposal.AlertMessage{
			ProjectName:  r.config.Config.App.Name,
			Environment:  r.config.Config.Environment,
			TraceID:      ctx.GetString(pkgtypes.TraceIdName),
			HOST:         ctx.Request.Host,
			URI:          ctx.Request.URL.RequestURI(),
			Method:       ctx.Request.Method,
			ErrorMessage: err,
			ErrorStack:   stack,
			Timestamp:    time.Now(),
		})
	}

	// 打印堆栈信息（如果配置允许）
	if r.config.PrintStack {
		fmt.Printf("Panic recovered: %v\nStack trace:\n%s\n", err, stack)
	}
}

// DefaultRecoveryConfig 返回默认 Recovery 配置。
//
// 参数:
//   - cfg: 应用配置
//   - alertNotify: 告警通知处理函数（可选）
//
// 返回:
//   - *RecoveryConfig: 默认 Recovery 配置
func DefaultRecoveryConfig(cfg *config.Config, alertNotify proposal.NotifyHandler) *RecoveryConfig {
	return &RecoveryConfig{
		Enabled:     true,
		AlertNotify: alertNotify,
		Config:      cfg,
		PrintStack:  true,
	}
}

// NewDefaultRecovery 创建默认 Recovery 中间件。
// 提供便捷的创建方式，使用默认配置。
//
// 参数:
//   - cfg: 应用配置
//   - logger: 日志记录器
//   - alertNotify: 告警通知处理函数（可选）
//
// 返回:
//   - *Recovery: Recovery 中间件实例
func NewDefaultRecovery(cfg *config.Config, logger *logger.Logger, alertNotify proposal.NotifyHandler) *Recovery {
	return NewRecovery(DefaultRecoveryConfig(cfg, alertNotify), logger)
}
