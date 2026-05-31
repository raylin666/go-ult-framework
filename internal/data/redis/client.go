// Package redis 提供 Redis 客户端封装。
package redis

import (
	"ult/pkg/repositories"

	"github.com/raylin666/go-utils/v2/cache/redis"
)

// NewDefaultClient 获取默认 Redis 连接的客户端实例。
//
// 参数:
//   - repo: Redis 仓库
//
// 返回:
//   - redis.Client: Redis 客户端
func NewDefaultClient(repo repositories.RedisRepo) redis.Client {
	return repo.Redis(repositories.RedisConnectionDefaultName).Get()
}
