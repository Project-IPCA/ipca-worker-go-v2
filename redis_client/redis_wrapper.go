package redis_client

import (
	"github.com/google/uuid"
)

type RedisMessage struct {
	Action string    `json:"action"`
	UserID uuid.UUID `json:"user_id"`
}

func (redisAction *RedisAction) NewMessage(action string, userId uuid.UUID) *RedisMessage {
	return &RedisMessage{
		Action: action,
		UserID: userId,
	}
}
