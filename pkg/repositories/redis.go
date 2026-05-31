// Package repositories 提供数据仓库抽象层。
package repositories

import (
	"ult/pkg/cache"
)

// Redis 连接默认名称常量。
const (
	RedisConnectionDefaultName = "default" // 默认 Redis 连接名称
)

var _ RedisRepo = (*redisRepo)(nil)

// RedisRepo Redis 仓库接口，定义 Redis 连接管理操作。
type RedisRepo interface {
	Count() int                   // 获取连接数量
	Has(name string) bool         // 检查连接是否存在
	Redis(name string) cache.Redis // 获取指定名称的 Redis 连接
	All() map[string]cache.Redis  // 获取所有 Redis 连接
}

// redisRepo Redis 仓库实例，管理多个 Redis 连接。
type redisRepo struct {
	resource map[string]cache.Redis // Redis 连接映射
}

// Count 获取 Redis 连接数量。
//
// 返回:
//   - int: 连接数量
func (repo *redisRepo) Count() int {
	return len(repo.resource)
}

// Has 检查指定名称的 Redis 连接是否存在。
//
// 参数:
//   - name: 连接名称
//
// 返回:
//   - bool: true 表示存在
func (repo *redisRepo) Has(name string) bool {
	if _, ok := repo.resource[name]; ok {
		return true
	}

	return false
}

// Redis 获取指定名称的 Redis 连接。
//
// 参数:
//   - name: 连接名称
//
// 返回:
//   - cache.Redis: Redis 连接实例
func (repo *redisRepo) Redis(name string) cache.Redis {
	return repo.resource[name]
}

// All 获取所有 Redis 连接。
//
// 返回:
//   - map[string]cache.Redis: Redis 连接映射
func (repo *redisRepo) All() map[string]cache.Redis {
	return repo.resource
}
