package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Project-IPCA/ipca-worker-go-v2/config"
	"github.com/Project-IPCA/ipca-worker-go-v2/db"
	"github.com/Project-IPCA/ipca-worker-go-v2/models"
	"github.com/Project-IPCA/ipca-worker-go-v2/rabbitmq_client"
	"github.com/Project-IPCA/ipca-worker-go-v2/redis_client"
	"github.com/Project-IPCA/ipca-worker-go-v2/service"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"C"
)

func main() {
	cfg := config.NewConfig()

	db_pool := db.Init(cfg)
	pubsub := redis_client.RedisClient(cfg)

	for {
		if err := pythonConsumer(db_pool, pubsub, cfg); err != nil {
			fmt.Printf("Consumer error: %v, retrying in 5 seconds...\n", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func pythonConsumer(db_pool *gorm.DB, pubsub *redis.Client, cfg *config.Config) error {
	rabbit, err := rabbitmq_client.RabbitMQClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer rabbit.Close()

	ch, err := rabbit.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		cfg.RabbitMQ.QueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	msgs, err := ch.Consume(
		cfg.RabbitMQ.QueueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %v", err)
	}

	connCloseChan := make(chan *amqp.Error)
	channelCloseChan := make(chan *amqp.Error)

	ch.NotifyClose(channelCloseChan)
	rabbit.NotifyClose(connCloseChan)

	fmt.Println("Connected to RabbitMQ. Waiting for messages...")

	for {
		select {
		case err := <-connCloseChan:
			return fmt.Errorf("connection closed: %v", err)

		case err := <-channelCloseChan:
			return fmt.Errorf("channel closed: %v", err)

		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("message channel closed")
			}

			if err := processMessage(ch, db_pool, msg, pubsub); err != nil {
				fmt.Printf("Error processing message: %v\n", err)
			}
		}
	}
}

func processMessage(ch *amqp.Channel, db_pool *gorm.DB, msg amqp.Delivery, pubsub *redis.Client) error {
	var msgBody models.ReciveMessage
	if err := json.Unmarshal(msg.Body, &msgBody); err != nil {
		msg.Nack(false, true)
		return fmt.Errorf("failed to parse message: %v", err)
	}

	fmt.Printf("Processing message: %+v\n", msgBody)

	switch msgBody.JobType {
	case "upsert-testcase":
		service.AddAndUpdateTestCase(ch, db_pool, msg, msgBody, pubsub)
	case "exercise-submit":
		service.RunSubmission(ch, db_pool, msg, msgBody, pubsub)
	default:
		fmt.Printf("Unknown job type: %s\n", msgBody.JobType)
		msg.Ack(false)
	}
	return nil
}
