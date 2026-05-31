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

func NewHTTPServer(
	config *config.Config,
	logger *logger.Logger,
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
		[]http.ServerOption{
			http.WithServerNetwork(config.Server.Http.Network),
			http.WithServerAddress(addr),
		},
		// pkg_http.EnableAlertNoxtify(email.NotifyHandler(ctx, config.Notify, logger)),
		pkg_http.EnableCors(cors_domains),
		pkg_http.EnablePProf())

	httpRouter(server)

	return server
}
