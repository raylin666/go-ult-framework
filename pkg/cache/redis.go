// Package cache 提供 Redis 缓存连接封装。
// 基于 go-utils/redis 实现，支持连接池管理和配置化连接。
package cache

import (
	"context"
	"fmt"
	"time"
	"ult/config/autoload"

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
}

// NewRedis 创建新的 Redis 连接实例。
// 根据配置初始化 Redis 客户端，支持连接池、超时等配置。
//
// 参数:
//   - name: Redis 连接名称，用于错误标识
//   - config: Redis 配置（地址、端口、密码、连接池等）
//
// 返回:
//   - Redis: Redis 客户端实例
//   - error: 连接创建错误
func NewRedis(name string, config autoload.Redis) (Redis, error) {
	var rds = new(redis)
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

	client, err := utils_redis.NewClient(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("new redis to %s client err", name)
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
