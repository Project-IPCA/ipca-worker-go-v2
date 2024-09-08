package config

import (
	"os"
)

type RabbitMQConfig struct {
	Host string
	Port string
	User string
	Password string
	QueueName string
}

func LoadRabbitMQConfig() RabbitMQConfig {
	return RabbitMQConfig{
		Host: os.Getenv("RABBITMQ_HOST"),
		Port: os.Getenv("RABBITMQ_PORT"),
		User: os.Getenv("RABBITMQ_USER"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
		QueueName: os.Getenv("RABBITMQ_QUEUENAME"),
	}
}
