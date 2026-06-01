// Package action 提供 Redis 分布式锁实现。
// 基于 Redis SET NX EX 实现安全的分布式锁功能。
package action

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/raylin666/go-utils/v2/cache/redis"
)

// Lock Redis 分布式锁。
type Lock struct {
	ctx      context.Context // 上下文
	client   redis.Client    // Redis 客户端
	key      string          // 锁键名
	value    string          // 锁值（UUID，用于标识锁持有者）
	duration time.Duration   // 锁过期时间
}

// NewLock 创建新的分布式锁实例。
// 默认过期时间为 1 秒。
// 使用 UUID 作为锁值，确保只有锁的持有者才能释放锁。
//
// 参数:
//   - ctx: 上下文
//   - client: Redis 客户端
//   - key: 锁键名
//
// 返回:
//   - *Lock: 分布式锁实例
func NewLock(ctx context.Context, client redis.Client, key string) *Lock {
	return &Lock{
		ctx:      ctx,
		client:   client,
		key:      key,
		value:    uuid.New().String(),
		duration: time.Second,
	}
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
// 使用 Redis SET NX EX 命令原子设置锁和过期时间。
// 锁的值为 UUID，用于标识锁的持有者。
//
// 返回:
//   - bool: true 表示成功获取锁
func (lock *Lock) Lock() bool {
	isOk, err := lock.client.SetNX(lock.ctx, lock.key, lock.value, lock.duration).Result()
	if err != nil {
		return false
	}

	return isOk
}

// UnLock 释放锁。
// 使用 Lua 脚本原子检查锁的值并删除。
// 只有锁的持有者（value 匹配）才能释放锁，防止误释放他人的锁。
func (lock *Lock) UnLock() {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`

	_, _ = lock.client.Eval(lock.ctx, script, []string{lock.key}, lock.value).Result()
}

// TryLock 尝试获取锁，支持重试。
// 在指定时间内多次尝试获取锁，直到成功或超时。
//
// 参数:
//   - retryInterval: 重试间隔
//   - timeout: 总超时时间
//
// 返回:
//   - bool: true 表示成功获取锁
func (lock *Lock) TryLock(retryInterval, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(lock.ctx, timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return false
		default:
			if lock.Lock() {
				return true
			}
			time.Sleep(retryInterval)
		}
	}
}

// IsLocked 检查锁是否被持有（不检查持有者）。
//
// 返回:
//   - bool: true 表示锁被持有
func (lock *Lock) IsLocked() bool {
	value, err := lock.client.Get(lock.ctx, lock.key).Result()
	if err != nil {
		return false
	}
	return value != ""
}

// IsHeldByMe 检查锁是否由当前实例持有。
//
// 返回:
//   - bool: true 表示锁由当前实例持有
func (lock *Lock) IsHeldByMe() bool {
	value, err := lock.client.Get(lock.ctx, lock.key).Result()
	if err != nil {
		return false
	}
	return value == lock.value
}
