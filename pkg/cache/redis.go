package cache

import (
	"context"
	"fmt"
	"time"
	"ult/config/autoload"

	utils_redis "github.com/raylin666/go-utils/v2/cache/redis"
)

var _ Redis = (*redis)(nil)

type Redis interface {
	Get() utils_redis.Client
	Close() error
}

type redis struct {
	client utils_redis.Client
}

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

func (rds *redis) Get() utils_redis.Client {
	return rds.client
}

func (rds *redis) Close() error {
	return rds.Get().Close()
}
