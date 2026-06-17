// Package server 提供服务器层实现。
// 服务器层负责创建和配置 HTTP 服务器实例。
package server

import (
	"fmt"
	"strings"

	"ult/config"
	"ult/internal/router"
	pkg_http "ult/pkg/http"
	"ult/pkg/logger"

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
//   - *pkg_http.HTTPServer: HTTP 服务器实例
func NewHTTPServer(
	config *config.Config,
	logger *logger.Logger,
	httpRouter router.HTTPRouter) *pkg_http.HTTPServer {
	var addr = fmt.Sprintf("%s:%d", config.Server.Http.Host, config.Server.Http.Port)
	var corsDomains []string
	if config.Server.Http.Cors.Domains == "all" {
		corsDomains = append(corsDomains, "*")
	} else {
		corsDomains = strings.Split(config.Server.Http.Cors.Domains, ",")
	}

	var server = pkg_http.NewServer(
		config,
		logger,
		[]http.ServerOption{
			http.WithServerNetwork(config.Server.Http.Network),
			http.WithServerAddress(addr),
		},
		// pkg_http.EnableAlertNotify(email.NotifyHandler(ctx, config.Notify, logger)),
		pkg_http.EnableCors(corsDomains),
		pkg_http.EnablePProf())

	// 注册路由
	httpRouter(server)

	return server
}

// Note: 新的中间件系统已集成到 pkg/http/server.go 中
// 默认中间件（CORS、Recovery、Request、Response）会自动注册
// 如需添加自定义中间件，可使用 server.UseMiddleware() 方法
