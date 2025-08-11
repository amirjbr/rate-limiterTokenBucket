package redisClient

import (
	"TokenBucketRateLimiter/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	//addr := fmt.Sprintf("%v:%v", cfg.Redis.Host, cfg.Redis.Port)
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	//Addr:     addr,
	//Password: cfg.Redis.Password})
}
