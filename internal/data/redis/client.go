package redis

import (
	"github.com/raylin666/go-utils/v2/cache/redis"
	"ult/pkg/repositories"
)

func NewDefaultClient(repo repositories.RedisRepo) redis.Client {
	return repo.Redis(repositories.RedisConnectionDefaultName).Get()
}