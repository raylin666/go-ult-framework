package server

import (
	"fmt"
	"github.com/raylin666/go-utils/http"
	"strings"
	"ult/config"
	"ult/internal/router"
	"ult/pkg/global"
	pkg_http "ult/pkg/http"
	"ult/pkg/logger"
)

// NewHTTPServer 创建 HTTP 服务
func NewHTTPServer(
	config *config.Config,
	logger *logger.Logger,
	dataRepo global.DataRepo,
	httpRouter router.HTTPRouter) *pkg_http.HTTPServer {
	var addr = fmt.Sprintf("%s:%d", config.Server.Http.Host, config.Server.Http.Port)
	var cors_domains []string
	if config.Server.Http.Cors.Domains == "all" {
		cors_domains = append(cors_domains, "*")
	} else {
		cors_domains = strings.Split(config.Server.Http.Cors.Domains, ",")
	}

	var server = pkg_http.NewServer(
		config,
		logger,
		dataRepo,
		[]http.ServerOption{
			http.WithServerNetwork(config.Server.Http.Network),
			http.WithServerAddress(addr),
		},
		pkg_http.EnableCors(cors_domains),
		pkg_http.EnablePProf())

	// 注册 HTTP 路由
	httpRouter(server)

	return server
}
