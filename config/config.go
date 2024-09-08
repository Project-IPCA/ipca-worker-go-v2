package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   DBConfig
	Redis RedisConfig
	RabbitMQ RabbitMQConfig
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	return &Config{
		DB:   LoadDBConfig(),
		Redis : LoadRedisConfig(),
		RabbitMQ: LoadRabbitMQConfig(),
	}
}
