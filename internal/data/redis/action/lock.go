package action

import (
	"context"
	"time"

	"github.com/raylin666/go-utils/v2/cache/redis"
)

type Lock struct {
	ctx      context.Context
	client   redis.Client
	key      string
	duration time.Duration
}

func NewLock(ctx context.Context, client redis.Client, key string) *Lock {
	return &Lock{ctx, client, key, time.Second}
}

func (lock *Lock) WithDuration(duration time.Duration) *Lock {
	lock.duration = duration
	return lock
}

func (lock *Lock) Lock() bool {
	var value = 1
	isOk, err := lock.client.SetNX(lock.ctx, lock.key, value, lock.duration).Result()
	if err != nil {
		return false
	}

	return isOk
}

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