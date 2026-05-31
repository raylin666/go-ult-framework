package main

import (
	"ult/config"
	genDb "ult/generate/gormgen/db"
	"ult/pkg/db"
	"ult/pkg/logger"
	"ult/pkg/repositories"

	utils_logger "github.com/raylin666/go-utils/v2/logger"
	"github.com/raylin666/go-utils/v2/server/system"
)

func main() {
	// 初始化配置
	conf, err := config.New()
	if err != nil {
		panic(err)
	}

	// 初始化 Environment
	var environment = system.NewEnvironment(conf.Environment)
	conf.Environment = environment.Value()

	// 初始化 Datetime
	var datetime = system.NewDatetime(
		system.WithLocation(conf.Datetime.Location),
		system.WithCSTLayout(conf.Datetime.CSTLayout))

	log, err := logger.NewJSONLogger(
		utils_logger.WithField(utils_logger.AppKey, conf.App.Name),
		utils_logger.WithField(utils_logger.EnvironmentKey, conf.Environment),
		utils_logger.WithTimeLayout(datetime.CSTLayout()))
	if err != nil {
		panic(err)
	}

	defaultDB, err := db.NewDb(repositories.DbConnectionDefaultName, conf.DB[repositories.DbConnectionDefaultName], log)
	if err != nil {
		panic(err)
	}

	// 生成文件存放目录
	var outPath = "../../internal/data/dbquery"

	// 执行生成默认数据库对应的模型文件
	genDb.NewGeneratorDefaultDb(defaultDB, outPath)
}
