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
	"ult/pkg/app"
	"ult/pkg/logger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/raylin666/go-utils/v2/http"
	"github.com/raylin666/go-utils/v2/middleware"
	utilsserver "github.com/raylin666/go-utils/v2/server"
	"github.com/raylin666/go-utils/v2/server/system"
)

// HTTPServer HTTP 服务器实例，实现了 app.Server 接口。
// 封装了 Gin 引擎，提供中间件链、CORS 处理、Panic 恢复和结构化日志等功能。
var _ app.Server = (*HTTPServer)(nil)

// HTTPServer 表示 HTTP 服务器实例。
// 包含 Gin 引擎、配置、日志记录器和服务器选项。
type HTTPServer struct {
	*option

	server *http.Server
	engine *gin.Engine
	config *config.Config
	logger *logger.Logger
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

	// 中间件处理
	srv.server.Handler = middleware.HTTPChain(srv.middlewares...)(srv)

	// 启动服务时自动打开浏览器操作
	if srv.option.openBrowser != "" {
		_ = system.OpenBrowser(srv.option.openBrowser)
	}

	// pprof 性能分析, register pprof to gin. 访问路径: /debug/pprof
	if srv.option.pprof && !system.NewEnvironment(srv.config.Environment).IsProd() {
		pprof.Register(srv.engine)
	}

	// 注册处理中间件
	srv.handlerMiddlewares()
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
