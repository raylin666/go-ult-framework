package repositories

import (
	"context"
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
	DbRepo() DbRepo
	Redis(name string) cache.Redis
	RedisRepo() RedisRepo
}

type dataRepo struct {
	db    *dbRepo
	redis *redisRepo
}

func NewDataRepo(logger *logger.Logger, db_config map[string]autoload.DB, redis_config map[string]autoload.Redis) DataRepo {
	var (
		ctx = context.Background()

		dbRepo    = new(dbRepo)
		redisRepo = new(redisRepo)
		repo      = new(dataRepo)
	)

	// 初始化数据库
	dbRepo.resource = make(map[string]db.Db, len(db_config))
	for dbName, dbConfig := range db_config {
		rdb, err := db.NewDb(dbName, dbConfig, logger)
		if err != nil {
			logger.UseApp(ctx).Error(fmt.Sprintf("init db.repo %s error", dbName), zap.Error(err))
		} else {
			logger.UseApp(ctx).Info(fmt.Sprintf("init db.repo %s success", dbName))
			dbRepo.resource[dbName] = rdb
		}
	}

	// 初始化 Redis
	redisRepo.resource = make(map[string]cache.Redis, len(redis_config))
	for redisName, redisConfig := range redis_config {
		redis, err := cache.NewRedis(redisName, redisConfig)
		if err != nil {
			logger.UseApp(ctx).Error(fmt.Sprintf("init redis.repo %s error", redisName), zap.Error(err))
		} else {
			logger.UseApp(ctx).Info(fmt.Sprintf("init redis.repo %s success", redisName))
			redisRepo.resource[redisName] = redis
		}
	}

	repo.db = dbRepo
	repo.redis = redisRepo
	return repo
}

func (repo *dataRepo) DB(name string) db.Db {
	return repo.db.resource[name]
}

func (repo *dataRepo) DbRepo() DbRepo {
	return repo.db
}

func (repo *dataRepo) Redis(name string) cache.Redis {
	return repo.redis.resource[name]
}

func (repo *dataRepo) RedisRepo() RedisRepo {
	return repo.redis
}
