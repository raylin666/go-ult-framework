package global

import (
	stdCtx "context"
	"fmt"
	"github.com/raylin666/go-utils/server/system"
	"go.uber.org/zap"
	"ult/config"
	"ult/pkg/logger"
)

var _ AppInterface = (*App)(nil)

type AppInterface interface {
	ID() string
	Name() string
	Version() string
	Environment() string
	Context() stdCtx.Context
	Run() error
}

type App struct {
	servers  []Server
	context  stdCtx.Context
	cancel   []func()
	shutdown system.Shutdown
	config   *config.Config
	logger   *logger.Logger
}

type AppOption func(*App)

func WithAppContext(ctx stdCtx.Context) AppOption {
	return func(app *App) {
		app.context = ctx
	}
}

func WithAppCancel(fn func()) AppOption {
	return func(app *App) {
		app.cancel = append(app.cancel, fn)
	}
}

func WithAppServer(srv ...Server) AppOption {
	return func(app *App) {
		app.servers = srv
	}
}

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

func (app *App) ID() string {
	return app.config.App.ID
}

func (app *App) Name() string {
	return app.config.App.Name
}

func (app *App) Version() string {
	return app.config.App.Version
}

func (app *App) Environment() string {
	return app.config.Environment
}

func (app *App) Context() stdCtx.Context {
	return app.context
}

func (app *App) Run() error {
	ctx := NewAppContext(app.context, app)
	for _, server := range app.servers {
		srvType := server.ServerType()
		// 注册服务关闭事件
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

		// 服务启动
		go func() {
			if err := server.Start(ctx); err != nil {
				app.logger.UseApp(ctx).Error(fmt.Sprintf("%s server startup err", srvType), zap.Error(err))
			}
		}()

		server.StartAfter()
	}

	// 处理服务关闭事件
	app.shutdown.Close(app.cancel...)
	return nil
}

type appKey struct{}

// NewAppContext returns a new Context that carries value.
func NewAppContext(ctx stdCtx.Context, s AppInterface) stdCtx.Context {
	return stdCtx.WithValue(ctx, appKey{}, s)
}

// FromAppContext returns the Transport value stored in ctx, if any.
func FromAppContext(ctx stdCtx.Context) (s AppInterface, ok bool) {
	s, ok = ctx.Value(appKey{}).(AppInterface)
	return
}
