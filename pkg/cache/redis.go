// Package cache 提供 Redis 缓存连接封装。
// 基于 go-utils/redis 实现，支持连接池管理、连接重试和配置化连接。
package cache

import (
	"context"
	"fmt"
	"time"
	"ult/config/autoload"
	"ult/pkg/logger"

	utils_redis "github.com/raylin666/go-utils/v2/cache/redis"
)

var _ Redis = (*redis)(nil)

// Redis Redis 客户端接口，定义缓存连接的基本操作。
type Redis interface {
	Get() utils_redis.Client        // 获取 Redis 客户端
	Close() error                   // 关闭 Redis 连接
	Ping(ctx context.Context) error // 测试 Redis 连接
}

// redis Redis 客户端实例，封装 go-utils Redis 客户端。
type redis struct {
	client utils_redis.Client
	logger *logger.Logger
}

// NewRedis 创建新的 Redis 连接实例。
// 支持连接重试机制，可配置最大重试次数和重试间隔。
//
// 参数:
//   - name: Redis 连接名称，用于日志标识
//   - config: Redis 配置（地址、端口、密码、连接池等）
//   - logger: 日志记录器实例
//
// 返回:
//   - Redis: Redis 客户端实例
//   - error: 连接创建错误
func NewRedis(name string, config autoload.Redis, logger *logger.Logger) (Redis, error) {
	var rds = new(redis)
	rds.logger = logger

	maxRetries := config.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	retryDelay := config.RetryDelay
	if retryDelay <= 0 {
		retryDelay = 1
	}

	opts := new(utils_redis.Options)
	opts.Addr = fmt.Sprintf("%s:%d", config.Addr, config.Port)
	opts.Network = config.Network
	opts.Username = config.Username
	opts.Password = config.Password
	opts.DB = config.DB
	opts.DialTimeout = time.Duration(config.DialTimeout)
	opts.ConnMaxIdleTime = time.Duration(config.IdleTimeout)
	opts.ConnMaxLifetime = time.Duration(config.MaxConnAge)
	opts.MaxRetries = config.MaxRetries
	opts.MaxRetryBackoff = time.Duration(config.MinRetryBackoff)
	opts.MinRetryBackoff = time.Duration(config.MinRetryBackoff)
	opts.MinIdleConns = config.MinIdleConns
	opts.WriteTimeout = time.Duration(config.WriteTimeout)
	opts.ReadTimeout = time.Duration(config.ReadTimeout)
	opts.PoolFIFO = config.PoolFIFO
	opts.PoolSize = config.PoolSize
	opts.PoolTimeout = time.Duration(config.PoolTimeout)

	var client utils_redis.Client
	var err error

	for i := 0; i < maxRetries; i++ {
		client, err = utils_redis.NewClient(context.TODO(), opts)
		if err == nil {
			break
		}

		logger.UseApp(context.TODO()).Warn(fmt.Sprintf("redis connection %s attempt %d/%d failed: %v, retrying in %d seconds", name, i+1, maxRetries, err, retryDelay))
		time.Sleep(time.Duration(retryDelay) * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("new redis to %s client err after %d retries", name, maxRetries)
	}

	rds.client = client

	return rds, nil
}

// Get 获取 Redis 客户端实例。
//
// 返回:
//   - utils_redis.Client: Redis 客户端
func (rds *redis) Get() utils_redis.Client {
	return rds.client
}

// Close 关闭 Redis 连接。
//
// 返回:
//   - error: 关闭连接时的错误
func (rds *redis) Close() error {
	return rds.Get().Close()
}

// Ping 测试 Redis 连接是否正常。
//
// 参数:
//   - ctx: 上下文
//
// 返回:
//   - error: 连接测试错误，nil 表示连接正常
func (rds *redis) Ping(ctx context.Context) error {
	return rds.Get().Ping(ctx).Err()
}
