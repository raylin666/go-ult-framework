package http

import (
	stdCtx "context"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/raylin666/go-utils/http"
	"github.com/raylin666/go-utils/middleware"
	"github.com/raylin666/go-utils/server"
	"github.com/raylin666/go-utils/server/system"
	nethttp "net/http"
	"time"
	"ult/config"
	"ult/pkg/global"
	"ult/pkg/logger"
	"ult/pkg/repositories"
)

var _ global.Server = (*HTTPServer)(nil)

type HTTPServer struct {
	*option

	server   *http.Server
	engine   *gin.Engine
	config   *config.Config
	logger   *logger.Logger
	dataRepo repositories.DataRepo
}

func NewServer(config *config.Config, log *logger.Logger, dataRepo repositories.DataRepo, srvOpts []http.ServerOption, opts ...Option) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	var srv = &HTTPServer{
		server:   http.NewServer(&nethttp.Server{}, srvOpts...),
		engine:   engine,
		config:   config,
		logger:   log,
		dataRepo: dataRepo,
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
	if srv.option.pprof && !srv.config.Env.IsProd() {
		pprof.Register(srv.engine)
	}

	// 注册处理中间件
	srv.handlerMiddlewares()
	return srv
}

func (srv *HTTPServer) Config() *config.Config {
	return srv.config
}

func (srv *HTTPServer) Logger() *logger.Logger {
	return srv.logger
}

func (srv *HTTPServer) DataRepo() repositories.DataRepo {
	return srv.dataRepo
}

// CreateRouterGroup 创建路由组
func (srv *HTTPServer) CreateRouterGroup() RouterGroup {
	return NewRouter(&srv.engine.RouterGroup)
}

func (srv *HTTPServer) ServeHTTP(writer nethttp.ResponseWriter, request *nethttp.Request) {
	srv.engine.ServeHTTP(writer, request)
}

// ServerAgreement 获取服务协议
func (srv *HTTPServer) ServerAgreement() *global.ServerAgreement {
	var agreement = new(global.ServerAgreement)
	agreement.Network = srv.config.Server.Http.Network
	agreement.Addr = fmt.Sprintf("%s:%d", srv.config.Server.Http.Host, srv.config.Server.Http.Port)
	agreement.Target = fmt.Sprintf("%s://%s", agreement.Network, agreement.Addr)
	return agreement
}

func (srv *HTTPServer) ServerType() string {
	return fmt.Sprintf("%s [%s]", server.HTTPServerType, srv.ServerAgreement().Addr)
}

func (srv *HTTPServer) StartBefore() {}

func (srv *HTTPServer) StartAfter() {}

func (srv *HTTPServer) CancelBefore() {}

func (srv *HTTPServer) CancelAfter() {}

func (srv *HTTPServer) Start(ctx stdCtx.Context) error {
	srv.logger.UseApp().Info(fmt.Sprintf("Serving HTTP-Server on %s", srv.ServerAgreement().Target))
	return srv.server.Start(ctx)
}

func (srv *HTTPServer) Stop(ctx stdCtx.Context) error {
	ctx, cancel := stdCtx.WithTimeout(ctx, srv.option.timeout)
	defer cancel()
	return srv.server.Stop(ctx)
}
