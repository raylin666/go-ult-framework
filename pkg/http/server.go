// Package http 提供 HTTP 服务器实现，基于 Gin 框架封装。
// 该包封装了常用的 HTTP 服务器功能，包括请求处理、中间件管理、CORS 支持和优雅关闭。
// 实现了自定义的 Context 包装器，提供请求验证、错误处理和响应格式化等额外功能。
package http

import (
	stdCtx "context"
	"fmt"
	nethttp "net/http"
	"time"
	"ult/config"
	"ult/errcode"
	"ult/pkg/app"
	pkgmiddleware "ult/pkg/http/middleware"
	"ult/pkg/logger"
	"ult/pkg/proposal"
	pkgtypes "ult/pkg/types"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/raylin666/go-utils/v2/http"
	utilsserver "github.com/raylin666/go-utils/v2/server"
	"github.com/raylin666/go-utils/v2/server/system"
	"github.com/raylin666/go-utils/v2/validator"
)

// HTTPServer HTTP 服务器实例，实现了 app.Server 接口。
// 封装了 Gin 引擎，提供中间件链、CORS 处理、Panic 恢复和结构化日志等功能。
var _ app.Server = (*HTTPServer)(nil)

// HTTPServer 表示 HTTP 服务器实例。
// 包含 Gin 引擎、配置、日志记录器和服务器选项。
type HTTPServer struct {
	*option

	server     *http.Server
	engine     *gin.Engine
	config     *config.Config
	logger     *logger.Logger
	middleware *pkgmiddleware.Manager // 中间件管理器
	validator  validator.Validator    // 数据验证器
}

// NewServer 创建新的 HTTPServer 实例。
// 初始化 Gin 引擎（release 模式），设置中间件链，可选启用 pprof 性能分析（非生产环境）。
//
// 参数:
//   - config: 应用配置
//   - log: 日志记录器实例
//   - srvOpts: go-utils HTTP 服务器选项
//   - opts: 自定义 HTTPServer 选项（cors、pprof、timeout 等）
//
// 返回:
//   - *HTTPServer: 初始化后的 HTTP 服务器实例
func NewServer(config *config.Config, log *logger.Logger, srvOpts []http.ServerOption, opts ...Option) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	var srv = &HTTPServer{
		server: http.NewServer(&nethttp.Server{}, srvOpts...),
		engine: engine,
		config: config,
		logger: log,
	}

	srv.option = &option{timeout: 5 * time.Second}
	for _, opt := range opts {
		opt(srv.option)
	}

	// 初始化中间件管理器
	srv.middleware = pkgmiddleware.NewManager()

	// 初始化数据验证器
	srv.validator = validator.New(
		validator.WithLocale(srv.config.Validator.Locale),
		validator.WithTagname(srv.config.Validator.Tagname))

	// 注册默认中间件
	srv.registerDefaultMiddlewares()

	// 注册通过 WithMiddleware 添加的自定义中间件
	for _, m := range srv.option.middlewares {
		srv.middleware.Use(m)
	}

	// 构建中间件链并应用到 Gin 引擎
	handlers := srv.middleware.Build()
	for _, handler := range handlers {
		engine.Use(gin.HandlerFunc(handler))
	}

	// 注册请求处理中间件（初始化 Context、验证器、响应处理）
	srv.registerRequestHandler()

	// 设置服务器处理器
	srv.server.Handler = srv

	// 启动服务时自动打开浏览器操作
	if srv.option.openBrowser != "" {
		_ = system.OpenBrowser(srv.option.openBrowser)
	}

	// pprof 性能分析, register pprof to gin. 访问路径: /debug/pprof
	if srv.option.pprof && !system.NewEnvironment(srv.config.Environment).IsProd() {
		pprof.Register(srv.engine)
	}

	return srv
}

// Config 返回应用配置。
func (srv *HTTPServer) Config() *config.Config {
	return srv.config
}

// Logger 返回服务器关联的日志记录器实例。
func (srv *HTTPServer) Logger() *logger.Logger {
	return srv.logger
}

// CreateRouterGroup 从 Gin 引擎的根路由组创建新的路由组。
// 这是应用中注册路由的入口点。
func (srv *HTTPServer) CreateRouterGroup() RouterGroup {
	return NewRouter(&srv.engine.RouterGroup)
}

// ServeHTTP 实现 http.Handler 接口。
// 将请求处理委托给底层的 Gin 引擎。
func (srv *HTTPServer) ServeHTTP(writer nethttp.ResponseWriter, request *nethttp.Request) {
	srv.engine.ServeHTTP(writer, request)
}

// ServerAgreement 返回服务器网络协议和地址信息。
// 包含网络类型（tcp）、主机地址和完整目标 URL。
func (srv *HTTPServer) ServerAgreement() *app.ServerAgreement {
	var agreement = new(app.ServerAgreement)
	agreement.Network = srv.config.Server.Http.Network
	agreement.Addr = fmt.Sprintf("%s:%d", srv.config.Server.Http.Host, srv.config.Server.Http.Port)
	agreement.Target = fmt.Sprintf("%s://%s", agreement.Network, agreement.Addr)
	return agreement
}

// ServerType 返回服务器类型和地址的描述字符串。
func (srv *HTTPServer) ServerType() string {
	return fmt.Sprintf("%s [%s]", utilsserver.HTTPServerType, srv.ServerAgreement().Addr)
}

// StartBefore 服务器启动前的钩子方法。
// 可被重写以执行启动前初始化操作。
func (srv *HTTPServer) StartBefore() {}

// StartAfter 服务器启动后的钩子方法。
// 可被重写以执行启动后操作。
func (srv *HTTPServer) StartAfter() {}

// CancelBefore 服务器停止前的钩子方法。
// 可被重写以执行停止前清理操作。
func (srv *HTTPServer) CancelBefore() {}

// CancelAfter 服务器停止后的钩子方法。
// 可被重写以执行停止后清理操作。
func (srv *HTTPServer) CancelAfter() {}

// Start 开始监听并提供 HTTP 请求服务。
// 记录服务器地址并启动底层 HTTP 服务器。
//
// 参数:
//   - ctx: 服务器生命周期的上下文
//
// 返回:
//   - error: 服务器启动过程中发生的任何错误
func (srv *HTTPServer) Start(ctx stdCtx.Context) error {
	srv.logger.UseApp(ctx).Info(fmt.Sprintf("Serving HTTP-Server on %s", srv.ServerAgreement().Target))
	return srv.server.Start(ctx)
}

// Stop 优雅关闭 HTTP 服务器，使用超时机制。
// 使用配置的超时时间进行优雅关闭。
//
// 参数:
//   - ctx: 关闭操作的上下文
//
// 返回:
//   - error: 关闭过程中发生的任何错误
func (srv *HTTPServer) Stop(ctx stdCtx.Context) error {
	ctx, cancel := stdCtx.WithTimeout(ctx, srv.option.timeout)
	defer cancel()
	return srv.server.Stop(ctx)
}

// registerDefaultMiddlewares 注册默认中间件。
// 包括 CORS 和 Recovery 中间件。
// Request 和 Response 处理逻辑保留在 server.go 中。
func (srv *HTTPServer) registerDefaultMiddlewares() {
	// 注册 CORS 中间件
	if len(srv.option.cors.domains) > 0 {
		srv.middleware.Use(pkgmiddleware.NewCORS(&pkgmiddleware.CORSConfig{
			Enabled:            true,
			AllowedOrigins:     srv.option.cors.domains,
			AllowCredentials:   true,
			OptionsPassthrough: true,
		}))
	}

	// 注册 Recovery 中间件
	srv.middleware.Use(pkgmiddleware.NewRecovery(
		&pkgmiddleware.RecoveryConfig{
			Enabled:     true,
			AlertNotify: srv.option.alertNotify,
			Config:      srv.config,
			PrintStack:  true,
		},
		srv.logger,
	))
}

// UseMiddleware 添加自定义中间件。
// 支持链式调用，可以连续添加多个中间件。
//
// 参数:
//   - m: 要添加的中间件实例
//
// 返回:
//   - *HTTPServer: HTTP 服务器实例（支持链式调用）
func (srv *HTTPServer) UseMiddleware(m pkgmiddleware.Middleware) *HTTPServer {
	srv.middleware.Use(m)
	return srv
}

// UseMiddlewareFunc 使用函数方式添加自定义中间件。
// 提供简化的中间件添加方式。
//
// 参数:
//   - name: 中间件名称
//   - priority: 中间件优先级
//   - handler: 中间件处理函数
//
// 返回:
//   - *HTTPServer: HTTP 服务器实例（支持链式调用）
func (srv *HTTPServer) UseMiddlewareFunc(name string, priority pkgmiddleware.Priority, handler pkgmiddleware.HandlerFunc) *HTTPServer {
	srv.middleware.UseFunc(name, priority, handler)
	return srv
}

// registerRequestHandler 注册请求处理函数。
// 初始化 Context、设置验证器、处理响应。
func (srv *HTTPServer) registerRequestHandler() {
	// 请求处理中间件
	srv.engine.Use(func(ctx *gin.Context) {
		// 拦截 404 请求路由
		if ctx.Writer.Status() == nethttp.StatusNotFound {
			return
		}

		// 请求时间
		ts := time.Now()

		// 初始化核心上下文 Context
		var appCtx = newContext(ctx)
		defer recoveryContext(appCtx)
		appCtx.init()
		appCtx.WithValidator(srv.validator)
		ctx.Set(pkgtypes.CoreContextNameKey, appCtx)

		defer func() {
			// 响应处理
			srv.handlerResponse(ts, ctx)
		}()

		ctx.Next()
	})
}

// handlerResponseHandler 响应处理。
// 根据请求状态生成成功或错误响应，记录请求日志。
//
// 参数:
//   - reqTime: 请求开始时间
//   - ctx: Gin 上下文
func (srv *HTTPServer) handlerResponse(reqTime time.Time, ctx *gin.Context) {
	var (
		resp            interface{}
		httpCode        int
		businessCode    int
		businessMessage string
		stackErr        error
		traceId         string
	)

	// 获取核心上下文 Context
	appCtx, ok := ctx.Value(pkgtypes.CoreContextNameKey).(Context)
	if !ok || appCtx == nil {
		return
	}

	// 获取链路追踪 TraceId
	traceId = appCtx.TraceID()

	// 发生错误, 进行返回
	if ctx.IsAborted() {
		if err := appCtx.GetAbortError(); err != nil {
			httpCode = err.HTTPCode()
			businessCode = err.BusinessCode()
			businessMessage = err.Message()
			stackErr = err.StackError()
			// 设置告警提醒 (如发邮件通知、如钉钉告警)
			if err.IsAlert() {
				if notifyHandler := srv.option.alertNotify; notifyHandler != nil {
					notifyHandler(&proposal.AlertMessage{
						ProjectName:  srv.config.App.Name,
						Environment:  srv.config.Environment,
						TraceID:      traceId,
						HOST:         appCtx.Host(),
						URI:          appCtx.URI(),
						Method:       appCtx.Method(),
						ErrorMessage: err,
						ErrorStack:   fmt.Sprintf("%+v", stackErr),
						Timestamp:    time.Now(),
					})
				}
			}

			resp = NewErrorResponse(traceId, businessCode, businessMessage, err.Desc())
		} else {
			err := errcode.ErrUnknownError
			httpCode = err.HTTPCode()
			businessCode = err.BusinessCode()
			businessMessage = err.Message()
			stackErr = ctx.Err()
			resp = NewErrorResponse(traceId, businessCode, businessMessage, err.Desc())
		}

		ctx.JSON(httpCode, resp)
	} else {
		// 响应正确返回
		httpCode = nethttp.StatusOK
		businessMessage = "OK"
		resp = NewSuccessResponse(traceId, appCtx.GetPayload())
		ctx.JSON(httpCode, resp)
	}

	costSeconds := time.Since(reqTime).Seconds()

	// 请求日志打印
	srv.logger.RequestLog(ctx, &logger.RequestLogFormat{
		ClientIp:          ctx.ClientIP(),
		Method:            appCtx.Method(),
		Path:              appCtx.URI(),
		RequestProto:      ctx.Request.Proto,
		RequestReferer:    ctx.Request.Referer(),
		RequestUa:         ctx.Request.UserAgent(),
		RequestPostData:   ctx.Request.PostForm.Encode(),
		RequestBodyData:   string(appCtx.RawData()),
		RequestHeaderData: appCtx.Header(),
		HttpCode:          ctx.Writer.Status(),
		BusinessCode:      businessCode,
		BusinessMessage:   businessMessage,
		RequestTime:       reqTime,
		ResponseTime:      time.Now(),
		CostSeconds:       costSeconds,
	}, stackErr)
}
