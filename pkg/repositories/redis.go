package repositories

import (
	"ult/pkg/cache"
)

const (
	RedisConnectionDefaultName = "default"
)

var _ RedisRepo = (*redisRepo)(nil)

type RedisRepo interface {
	Count() int
	Has(name string) bool
	Redis(name string) cache.Redis
	All() map[string]cache.Redis
}

type redisRepo struct {
	resource map[string]cache.Redis
}

func (repo *redisRepo) Count() int {
	return len(repo.resource)
}

func (repo *redisRepo) Has(name string) bool {
	if _, ok := repo.resource[name]; ok {
		return true
	}

	return false
}

func (repo *redisRepo) Redis(name string) cache.Redis {
	return repo.resource[name]
}

func (repo *redisRepo) All() map[string]cache.Redis {
	return repo.resource
}
