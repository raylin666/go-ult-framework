package main

import (
	"context"
	"fmt"
	"github.com/raylin666/go-utils/auth"
	utils_logger "github.com/raylin666/go-utils/logger"
	"github.com/raylin666/go-utils/server/system"
	"ult/config"
	"ult/internal/app"
	"ult/internal/constant/errcode"
	"ult/pkg/code"
	"ult/pkg/global"
	"ult/pkg/http"
	"ult/pkg/logger"
)

func newApp(
	config *config.Config,
	logger *logger.Logger,
	hs *http.HTTPServer) *global.App {
	var appCtx = context.Background()
	var appCancel = func() {}
	return global.NewApp(
		config,
		logger,
		global.WithAppContext(appCtx),
		global.WithAppServer(hs),
		global.WithAppCancel(appCancel))
}

func main() {
	// 初始化配置
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	if conf.Env.IsDev() || conf.Env.IsTest() {
		fmt.Println(fmt.Sprintf("Warning: %s cannot be found, or it is illegal. The default %s will be used.", conf.Env.Value(), conf.Env.Value()))
	}

	// 打印启动信息
	app.NewLogo(conf)

	// 初始化错误状态码
	code.New(conf.Language.Local)
	errcode.RegisterNewMerged()

	// 初始化 Datetime
	app.Datetime = system.NewDatetime(
		system.WithLocation(conf.Datetime.Location),
		system.WithCSTLayout(conf.Datetime.CSTLayout))

	// 初始化 JWT 鉴权认证
	app.JWT = auth.NewJWT(conf.JWT.App, conf.JWT.Key, conf.JWT.Secret)

	// 初始化 Logger
	applogger, err := logger.NewJSONLogger(
		// utils_logger.WithDisableConsole(),
		utils_logger.WithField(utils_logger.AppKey, conf.App.Name),
		utils_logger.WithField(utils_logger.EnvironmentKey, conf.Env.Value()),
		utils_logger.WithTimeLayout(app.Datetime.CSTLayout()),
		utils_logger.WithPathFileRotation(fmt.Sprintf("%s/runtime/logs/%s.log", conf.ProjectPath, conf.App.Name), utils_logger.PathFileRotationOption{
			MaxSize:    conf.Logger.MaxSize,
			MaxAge:     conf.Logger.MaxAge,
			MaxBackups: conf.Logger.MaxBackups,
			LocalTime:  conf.Logger.LocalTime,
			Compress:   conf.Logger.Compress,
		})) //	项目访问日志存放文件
	if err != nil {
		panic(err)
	}

	// 初始化数据仓库
	var repo = global.NewDataRepo(applogger, conf.DB, conf.Redis)

	// 初始化应用服务
	application, cleanup, err := initApp(conf, applogger, repo)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = applogger.Sync()

		cleanup()
	}()

	// start and wait for stop signal
	if err := application.Run(); err != nil {
		panic(err)
	}
}
