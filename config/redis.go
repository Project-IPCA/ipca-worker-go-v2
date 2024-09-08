package config

import (
	"os"
)

type RedisConfig struct {
	Host string
	Port string
	User string
	Password string
}

func LoadRedisConfig() RedisConfig {
	return RedisConfig{
		Host: os.Getenv("REDIS_HOST"),
		Port: os.Getenv("REDIS_PORT"),
		User: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
}
