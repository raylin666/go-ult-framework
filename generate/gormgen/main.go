package main

import (
	utils_logger "github.com/raylin666/go-utils/logger"
	"github.com/raylin666/go-utils/server/system"
	"ult/config"
	genDb "ult/generate/gormgen/db"
	"ult/internal/app"
	"ult/internal/constant/defined"
	"ult/pkg/db"
	"ult/pkg/logger"
)

func main() {
	// 初始化配置
	conf, err := config.New()
	if err != nil {
		panic(err)
	}

	// 初始化 Datetime
	app.Datetime = system.NewDatetime(
		system.WithLocation(conf.Datetime.Location),
		system.WithCSTLayout(conf.Datetime.CSTLayout))

	log, err := logger.NewJSONLogger(
		utils_logger.WithField(utils_logger.AppKey, conf.App.Name),
		utils_logger.WithField(utils_logger.EnvironmentKey, conf.Env.Value()),
		utils_logger.WithTimeLayout(app.Datetime.CSTLayout()))
	if err != nil {
		panic(err)
	}

	rdb, err := db.NewDb(defined.DB_CONNECTION_DEFAULT_NAME, conf.DB[defined.DB_CONNECTION_DEFAULT_NAME], log)
	if err != nil {
		panic(err)
	}

	// 生成文件存放目录
	var outPath = "../../internal/repositories/dbrepo/query"

	// 执行生成默认数据库对应的模型文件
	genDb.NewGeneratorDefaultDb(rdb, outPath)
}
