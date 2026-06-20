// Package server 提供服务器层实现。
// 服务器层负责创建和配置 HTTP 服务器实例。
package server

import (
	"fmt"
	"strings"

	"ult/config"
	"ult/internal/app"
	"ult/internal/router"
	pkghttp "ult/pkg/http"
	pkgmiddleware "ult/pkg/http/middleware"

	"github.com/raylin666/go-utils/v2/http"
)

// NewHTTPServer 创建并配置 HTTP 服务器实例。
// 设置服务器地址、CORS、pprof 等配置，注册路由。
//
// 参数:
//   - config: 应用配置
//   - logger: 日志实例
//   - httpRouter: 路由注册函数
//
// 返回:
//   - *pkghttp.HTTPServer: HTTP 服务器实例
func NewHTTPServer(
	config *config.Config,
	tool *app.Tools,
	httpRouter router.HTTPRouter) *pkghttp.HTTPServer {
	var addr = fmt.Sprintf("%s:%d", config.Server.Http.Host, config.Server.Http.Port)
	var corsDomains []string
	if config.Server.Http.Cors.Domains == "all" {
		corsDomains = append(corsDomains, "*")
	} else {
		corsDomains = strings.Split(config.Server.Http.Cors.Domains, ",")
	}

	// 创建 CORS 中间件
	corsMiddleware := pkgmiddleware.NewCORS(&pkgmiddleware.CORSConfig{
		Enabled:            len(corsDomains) > 0,
		AllowedOrigins:     corsDomains,
		AllowCredentials:   true,
		OptionsPassthrough: true,
	})

	// 创建 PProf 中间件
	pprofMiddleware := pkgmiddleware.NewPProf(&pkgmiddleware.PProfConfig{
		Enabled:     true,
		Environment: config.Environment,
	})

	// 创建 Recovery 中间件（默认启用）
	recoveryMiddleware := pkgmiddleware.NewDefaultRecovery(config, tool.Logger(), nil)

	// 创建服务器
	var server = pkghttp.NewServer(
		config,
		tool.Logger(),
		[]http.ServerOption{
			http.WithServerNetwork(config.Server.Http.Network),
			http.WithServerAddress(addr),
		},
		pkghttp.WithMiddleware(
			recoveryMiddleware,
			corsMiddleware,
			pprofMiddleware),
	)

	// 注册路由
	httpRouter(server)

	// 注册 PProf 路由（特殊处理）
	pprofMiddleware.RegisterRoutes(server.Engine())

	return server
}
