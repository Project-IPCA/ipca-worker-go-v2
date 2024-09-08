package redis_client

import (
	"fmt"

	"github.com/Project-IPCA/ipca-worker-go-v2/config"
	"github.com/redis/go-redis/v9"
)

func RedisClient(cfg *config.Config) *redis.Client {
	url := fmt.Sprintf("redis://%s:%s@%s:%s/", cfg.Redis.User, cfg.Redis.Password, cfg.Redis.Host, cfg.Redis.Port)

	opt, err := redis.ParseURL(url)
	if err != nil {
		panic("failed to connect to redis: " + err.Error())
	}

	return redis.NewClient(opt)
}
