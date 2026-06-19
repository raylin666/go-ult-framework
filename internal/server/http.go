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

	var server = pkghttp.NewServer(
		config,
		tool.Logger(),
		[]http.ServerOption{
			http.WithServerNetwork(config.Server.Http.Network),
			http.WithServerAddress(addr),
		},
		pkghttp.EnableCors(corsDomains),
		pkghttp.EnablePProf())

	// 注册路由
	httpRouter(server)

	return server
}