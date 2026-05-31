// Package app 提供应用管理功能。
// 定义应用生命周期管理、服务器注册和上下文管理等功能。
package app

import (
	stdCtx "context"
	"fmt"
	"ult/config"
	"ult/pkg/logger"
	"ult/pkg/types"

	"github.com/raylin666/go-utils/v2/server/system"
	"go.uber.org/zap"
)

var _ AppInterface = (*App)(nil)

// AppInterface 应用接口，定义应用的基本属性和操作。
type AppInterface interface {
	ID() string              // 获取应用 ID
	Name() string            // 获取应用名称
	Version() string         // 获取应用版本
	Environment() string     // 获取运行环境
	Context() stdCtx.Context // 获取应用上下文
	Run() error              // 运行应用
}

// App 应用实例，管理服务器生命周期和应用配置。
type App struct {
	servers  []Server        // 服务器列表
	context  stdCtx.Context  // 应用上下文
	cancel   []func()        // 取消函数列表
	shutdown system.Shutdown // 关闭处理器
	config   *config.Config  // 应用配置
	logger   *logger.Logger  // 日志记录器
}

// AppOption 应用选项函数类型。
type AppOption func(*App)

// WithAppContext 设置应用上下文选项。
//
// 参数:
//   - ctx: 应用上下文
//
// 返回:
//   - AppOption: 应用选项函数
func WithAppContext(ctx stdCtx.Context) AppOption {
	return func(app *App) {
		app.context = ctx
	}
}

// WithAppCancel 设置应用取消函数选项。
//
// 参数:
//   - fn: 取消函数
//
// 返回:
//   - AppOption: 应用选项函数
func WithAppCancel(fn func()) AppOption {
	return func(app *App) {
		app.cancel = append(app.cancel, fn)
	}
}

// WithAppServer 设置应用服务器选项。
//
// 参数:
//   - srv: 服务器列表
//
// 返回:
//   - AppOption: 应用选项函数
func WithAppServer(srv ...Server) AppOption {
	return func(app *App) {
		app.servers = srv
	}
}

// NewApp 创建新的应用实例。
// 初始化应用配置、日志记录器和关闭处理器。
//
// 参数:
//   - config: 应用配置
//   - logger: 日志记录器
//   - opts: 应用选项列表
//
// 返回:
//   - *App: 应用实例
func NewApp(config *config.Config, logger *logger.Logger, opts ...AppOption) *App {
	var app = &App{
		config:   config,
		logger:   logger,
		context:  stdCtx.Background(),
		shutdown: system.NewShutdown(),
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

// ID 获取应用 ID。
//
// 返回:
//   - string: 应用 ID
func (app *App) ID() string {
	return app.config.App.ID
}

// Name 获取应用名称。
//
// 返回:
//   - string: 应用名称
func (app *App) Name() string {
	return app.config.App.Name
}

// Version 获取应用版本。
//
// 返回:
//   - string: 应用版本
func (app *App) Version() string {
	return app.config.App.Version
}

// Environment 获取运行环境。
//
// 返回:
//   - string: 运行环境（dev、prod 等）
func (app *App) Environment() string {
	return app.config.Environment
}

// Context 获取应用上下文。
//
// 返回:
//   - stdCtx.Context: 应用上下文
func (app *App) Context() stdCtx.Context {
	return app.context
}

// Run 运行应用。
// 启动所有注册的服务器，并设置优雅关闭处理。
//
// 返回:
//   - error: 运行错误
func (app *App) Run() error {
	ctx := NewAppContext(app.context, app)
	for _, server := range app.servers {
		srvType := server.ServerType()
		app.cancel = append(app.cancel, func() {
			server.CancelBefore()
			if err := server.Stop(ctx); err != nil {
				app.logger.UseApp(ctx).Error(fmt.Sprintf("%s server shutdown err", srvType), zap.Error(err))
			} else {
				app.logger.UseApp(ctx).Info(fmt.Sprintf("%s server is success close", srvType))
			}

			server.CancelAfter()
		})

		server.StartBefore()

		go func() {
			if err := server.Start(ctx); err != nil {
				app.logger.UseApp(ctx).Error(fmt.Sprintf("%s server startup err", srvType), zap.Error(err))
			}
		}()

		server.StartAfter()
	}

	app.shutdown.Close(app.cancel...)
	return nil
}

// appKey 应用上下文键。
type appKey struct{}

// NewAppContext 创建带有应用实例的上下文。
//
// 参数:
//   - ctx: 基础上下文
//   - s: 应用实例
//
// 返回:
//   - stdCtx.Context: 带有应用实例的上下文
func NewAppContext(ctx stdCtx.Context, s AppInterface) stdCtx.Context {
	return stdCtx.WithValue(ctx, appKey{}, s)
}

// FromAppContext 从上下文中获取应用实例。
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - s: 应用实例
//   - ok: 是否存在
func FromAppContext(ctx stdCtx.Context) (s AppInterface, ok bool) {
	s, ok = ctx.Value(appKey{}).(AppInterface)
	return
}

// NewAppContextWithTraceId 创建带有应用实例和 TraceID 的上下文。
//
// 参数:
//   - ctx: 基础上下文
//   - app: 应用实例
//   - traceId: 链路追踪 ID
//
// 返回:
//   - stdCtx.Context: 带有应用实例和 TraceID 的上下文
func NewAppContextWithTraceId(ctx stdCtx.Context, app AppInterface, traceId string) stdCtx.Context {
	ctx = NewAppContext(ctx, app)
	var reqContext = new(types.RequestContext)
	reqContext.WithTraceID(traceId)
	return types.NewRequestContext(ctx, reqContext)
}
