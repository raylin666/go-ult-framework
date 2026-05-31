//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"ult/config"
	"ult/internal/api"
	"ult/internal/app"
	"ult/internal/data"
	"ult/internal/data/repo"
	"ult/internal/router"
	"ult/internal/server"
	"ult/internal/service"
	pkgapp "ult/pkg/app"
	"ult/pkg/logger"

	"github.com/google/wire"
)

func initApp(conf *config.Config, tools *app.Tools) (*pkgapp.App, func(), error) {
	panic(wire.Build(
		data.ProviderSet,
		repo.ProviderSet,
		service.ProviderSet,
		api.ProviderSet,
		router.ProviderSet,
		server.ProviderSet,
		provideLogger,
		newApp))
}

func provideLogger(tools *app.Tools) *logger.Logger {
	return tools.Logger()
}

