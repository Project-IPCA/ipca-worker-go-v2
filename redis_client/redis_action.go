package redis_client

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type IRedisAction interface {
	PublishMessage(channel, message string)
	SubscribeTopic(channel string)
}

type RedisAction struct {
	Redis *redis.Client
}

func NewRedisAction(redis *redis.Client) *RedisAction {
	return &RedisAction{Redis: redis}
}

func (redisAction *RedisAction) PublishMessage(channel string, message interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	body, err := json.Marshal(message)
	if err != nil {
		panic("failed to marshal message to JSON: " + err.Error())
	}

	err = redisAction.Redis.Publish(ctx, channel, body).Err()
	if err != nil {
		panic("failed to publish message: " + err.Error())
	}

	return nil
}
