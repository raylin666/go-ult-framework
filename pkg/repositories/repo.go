// Package repositories 提供数据仓库抽象层。
package repositories

import (
	"context"
	"fmt"
	"ult/config"
	"ult/pkg/cache"
	"ult/pkg/db"
	"ult/pkg/logger"

	"go.uber.org/zap"
)

var _ DataRepo = (*dataRepo)(nil)

// DataRepo 数据仓库接口，定义数据库和 Redis 连接管理操作。
type DataRepo interface {
	DB(name string) db.Db          // 获取指定名称的数据库连接
	DbRepo() DbRepo                // 获取数据库仓库
	Redis(name string) cache.Redis // 获取指定名称的 Redis 连接
	RedisRepo() RedisRepo          // 获取 Redis 仓库
}

// dataRepo 数据仓库实例，管理数据库和 Redis 连接。
type dataRepo struct {
	db    *dbRepo    // 数据库仓库
	redis *redisRepo // Redis 仓库
}

// NewDataRepo 创建新的数据仓库实例。
// 根据配置初始化数据库和 Redis 连接。
//
// 参数:
//   - logger: 日志记录器
//   - conf: 应用配置
//
// 返回:
//   - DataRepo: 数据仓库实例
func NewDataRepo(logger *logger.Logger, conf *config.Config) DataRepo {
	var (
		ctx = context.Background()

		dbRepo    = new(dbRepo)
		redisRepo = new(redisRepo)
		repo      = new(dataRepo)

		dbMap    = conf.DB
		redisMap = conf.Redis
	)

	// 初始化数据库
	dbMap[DbConnectionDefaultName] = conf.DB["default"]

	lenDatabase := len(dbMap)
	if lenDatabase > 0 {
		dbRepo.resource = make(map[string]db.Db, lenDatabase)
		for dbName, dbConfig := range dbMap {
			rdb, err := db.NewDb(dbName, dbConfig, logger)
			if err != nil {
				logger.UseApp(ctx).Error(fmt.Sprintf("初始化 Db.repo %s 失败", dbName), zap.Error(err))
			} else {
				logger.UseApp(ctx).Info(fmt.Sprintf("初始化 Db.repo %s 成功", dbName))
				dbRepo.resource[dbName] = rdb
			}
		}

		repo.db = dbRepo
	} else {
		logger.UseApp(ctx).Warn("目前 Db.repo 未配置连接")
	}

	// 初始化 Redis
	redisMap[RedisConnectionDefaultName] = conf.Redis["default"]

	lenRedis := len(redisMap)
	if lenRedis > 0 {
		redisRepo.resource = make(map[string]cache.Redis, lenRedis)
		for redisName, redisConfig := range redisMap {
			redis, err := cache.NewRedis(redisName, redisConfig, logger)
			if err != nil {
				logger.UseApp(ctx).Error(fmt.Sprintf("初始化 Redis.repo %s 失败", redisName), zap.Error(err))
			} else {
				logger.UseApp(ctx).Info(fmt.Sprintf("初始化 Redis.repo %s 成功", redisName))
				redisRepo.resource[redisName] = redis
			}
		}

		repo.redis = redisRepo
	} else {
		logger.UseApp(ctx).Warn("目前 Redis.repo 未配置连接")
	}

	return repo
}

// DB 获取指定名称的数据库连接。
//
// 参数:
//   - name: 连接名称
//
// 返回:
//   - db.Db: 数据库连接实例
func (repo *dataRepo) DB(name string) db.Db {
	if repo.db == nil {
		return nil
	}
	return repo.db.resource[name]
}

// DbRepo 获取数据库仓库。
//
// 返回:
//   - DbRepo: 数据库仓库实例
func (repo *dataRepo) DbRepo() DbRepo {
	return repo.db
}

// Redis 获取指定名称的 Redis 连接。
//
// 参数:
//   - name: 连接名称
//
// 返回:
//   - cache.Redis: Redis 连接实例
func (repo *dataRepo) Redis(name string) cache.Redis {
	if repo.redis == nil {
		return nil
	}
	return repo.redis.resource[name]
}

// RedisRepo 获取 Redis 仓库。
//
// 返回:
//   - RedisRepo: Redis 仓库实例
func (repo *dataRepo) RedisRepo() RedisRepo {
	return repo.redis
}
