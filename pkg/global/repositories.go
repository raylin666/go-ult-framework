package global

import (
	"fmt"
	"go.uber.org/zap"
	"ult/config/autoload"
	"ult/pkg/cache"
	"ult/pkg/db"
	"ult/pkg/logger"
)

var _ DataRepo = (*dataRepo)(nil)

type DataRepo interface {
	DB(name string) db.Db
	AllDB() map[string]db.Db
	Redis(name string) cache.Redis
	AllRedis() map[string]cache.Redis
}

type dataRepo struct {
	db      map[string]db.Db
	redis   map[string]cache.Redis
}

func NewDataRepo(logger *logger.Logger, db_config map[string]autoload.DB, redis_config map[string]autoload.Redis) DataRepo {
	var repo = new(dataRepo)

	// 初始化数据库
	repo.db = make(map[string]db.Db, len(db_config))
	for dbName, dbConfig := range db_config {
		rdb, err := db.NewDb(dbName, dbConfig, logger)
		if err != nil {
			logger.UseApp().Error(fmt.Sprintf("init db.repo %s error", dbName), zap.Error(err))
		} else {
			logger.UseApp().Info(fmt.Sprintf("init db.repo %s success", dbName))
			repo.db[dbName] = rdb
		}
	}

	// 初始化 Redis
	repo.redis = make(map[string]cache.Redis, len(redis_config))
	for redisName, redisConfig := range redis_config {
		redis, err := cache.NewRedis(redisName, redisConfig)
		if err != nil {
			logger.UseApp().Error(fmt.Sprintf("init redis.repo %s error", redisName), zap.Error(err))
		} else {
			logger.UseApp().Info(fmt.Sprintf("init redis.repo %s success", redisName))
			repo.redis[redisName] = redis
		}
	}

	return repo
}

func (repo *dataRepo) DB(name string) db.Db {
	return repo.db[name]
}

func (repo *dataRepo) AllDB() map[string]db.Db {
	return repo.db
}

func (repo *dataRepo) Redis(name string) cache.Redis {
	return repo.redis[name]
}

func (repo *dataRepo) AllRedis() map[string]cache.Redis {
	return repo.redis
}
