package main

import (
	"context"
	"fmt"
	"ult/config"
	"ult/errcode"
	"ult/internal/app"
	pkgapp "ult/pkg/app"
	"ult/pkg/http"
	"ult/pkg/logger"

	"github.com/raylin666/go-utils/v2/auth"
	utils_logger "github.com/raylin666/go-utils/v2/logger"
	"github.com/raylin666/go-utils/v2/server/system"
)

func newApp(
	config *config.Config,
	logger *logger.Logger,
	hs *http.HTTPServer) *pkgapp.App {
	appCtx, appCancel := context.WithCancel(context.Background())
	return pkgapp.NewApp(
		config,
		logger,
		pkgapp.WithAppContext(appCtx),
		pkgapp.WithAppServer(hs),
		pkgapp.WithAppCancel(appCancel))
}

func main() {
	// 初始化配置
	conf, err := config.New()
	if err != nil {
		panic(fmt.Sprintf("Config initialization failed: %v", err))
	}

	// 初始化 Environment
	var environment = system.NewEnvironment(conf.Environment)
	conf.Environment = environment.Value()

	// 打印启动信息
	app.NewLogo(conf)

	// 初始化错误状态码
	errcode.NewRegistry(conf.Language.Local)

	// 初始化 Datetime
	var datetime = system.NewDatetime(
		system.WithLocation(conf.Datetime.Location),
		system.WithCSTLayout(conf.Datetime.CSTLayout))

	// 初始化 Logger
	appLogger, err := logger.NewJSONLogger(
		// utils_logger.WithDisableConsole(),
		utils_logger.WithField(utils_logger.AppKey, conf.App.Name),
		utils_logger.WithField(utils_logger.EnvironmentKey, conf.Environment),
		utils_logger.WithTimeLayout(datetime.CSTLayout()),
		utils_logger.WithPathFileRotation(fmt.Sprintf("%s/runtime/logs/%s.log", conf.ProjectPath, conf.App.Name), utils_logger.PathFileRotationOption{
			MaxSize:    conf.Logger.MaxSize,
			MaxAge:     conf.Logger.MaxAge,
			MaxBackups: conf.Logger.MaxBackups,
			LocalTime:  conf.Logger.LocalTime,
			Compress:   conf.Logger.Compress,
		})) //	项目访问日志存放文件
	if err != nil {
		panic(fmt.Sprintf("Logger initialization failed: %v", err))
	}

	// 初始化 JWT 鉴权认证
	jwt, err := auth.NewJWT(conf.JWT.App, conf.JWT.Key, conf.JWT.Secret)
	if err != nil {
		appLogger.UseApp(context.Background()).Error(fmt.Sprintf("JWT initialization failed: %v", err))
		panic(err)
	}

	// 创建公共工具实例
	var tools = app.NewTools(appLogger, datetime, environment, jwt)

	// 初始化应用服务
	application, cleanup, err := initApp(conf, tools)
	if err != nil {
		appLogger.UseApp(context.Background()).Error(fmt.Sprintf("Application initialization failed: %v", err))
		panic(err)
	}

	defer func() {
		_ = appLogger.Sync()

		cleanup()
	}()

	// start and wait for stop signal
	if err := application.Run(); err != nil {
		appLogger.UseApp(context.Background()).Error(fmt.Sprintf("Application run failed: %v", err))
		panic(err)
	}
}
