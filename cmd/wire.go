//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"ult/config"
	"ult/internal/router"
	"ult/internal/server"
	"ult/pkg/global"
	"ult/pkg/logger"
)

// initApp init application.
func initApp(
	conf *config.Config,
	log *logger.Logger,
	repo global.DataRepo) (*global.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		router.ProviderSet,
		newApp))
}
