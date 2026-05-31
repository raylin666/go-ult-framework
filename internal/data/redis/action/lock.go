// Package action 提供 Redis 分布式锁实现。
// 基于 Redis SETNX 实现简单的分布式锁功能。
package action

import (
	"context"
	"time"

	"github.com/raylin666/go-utils/v2/cache/redis"
)

// Lock Redis 分布式锁。
type Lock struct {
	ctx      context.Context // 上下文
	client   redis.Client    // Redis 客户端
	key      string          // 锁键名
	duration time.Duration   // 锁过期时间
}

// NewLock 创建新的分布式锁实例。
// 默认过期时间为 1 秒。
//
// 参数:
//   - ctx: 上下文
//   - client: Redis 客户端
//   - key: 锁键名
//
// 返回:
//   - *Lock: 分布式锁实例
func NewLock(ctx context.Context, client redis.Client, key string) *Lock {
	return &Lock{ctx, client, key, time.Second}
}

// WithDuration 设置锁过期时间。
//
// 参数:
//   - duration: 过期时间
//
// 返回:
//   - *Lock: 分布式锁实例
func (lock *Lock) WithDuration(duration time.Duration) *Lock {
	lock.duration = duration
	return lock
}

// Lock 尝试获取锁。
// 使用 Redis SETNX 命令实现，返回是否成功获取锁。
//
// 返回:
//   - bool: true 表示成功获取锁
func (lock *Lock) Lock() bool {
	var value = 1
	isOk, err := lock.client.SetNX(lock.ctx, lock.key, value, lock.duration).Result()
	if err != nil {
		return false
	}

	return isOk
}

// UnLock 释放锁。
// 使用 Lua 脚本确保只有锁的持有者才能释放锁。
func (lock *Lock) UnLock() {
	value, err := lock.client.Get(lock.ctx, lock.key).Result()
	if err != nil {
		return
	}

	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
	`

	_, _ = lock.client.Eval(lock.ctx, script, []string{lock.key}, value).Result()
}