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
)

func main(){
	cfg := config.NewConfig()

	db_pool := db.Init(cfg)
	pubsub := redis_client.RedisClient(cfg)
	rabbit := rabbitmq_client.RabbitMQClient(cfg)

	pythonConsumer(db_pool,pubsub,rabbit,cfg)
}

func pythonConsumer(db_pool *gorm.DB, pubsub *redis.Client, rabbit *amqp.Connection,cfg *config.Config) {
	for {
		ch, err := rabbit.Channel()
		if err != nil {
			fmt.Printf("Failed to create a RabbitMQ channel: %v", err)
			time.Sleep(5 * time.Second)
			return
		}
		defer ch.Close()

		// Assert the queue
		_, err = ch.QueueDeclare(
			cfg.RabbitMQ.QueueName, // name
			true,          // durable
			false,         // delete when unused
			false,         // exclusive
			false,         // no-wait
			nil,           // arguments
		)
		if err != nil {
			fmt.Printf("Failed to declare the RabbitMQ queue: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = ch.Qos(1, 0, false)
		if err != nil {
			fmt.Printf("Failed to set QoS: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		msgs, err := ch.Consume(
			cfg.RabbitMQ.QueueName, // queue
			"",            // consumer
			false,         // auto-ack
			false,         // exclusive
			false,         // no-local
			false,         // no-wait
			nil,           // args
		)
		if err != nil {
			fmt.Printf("Failed to register RabbitMQ consumer: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println("Waiting for messages...")

		for msg := range msgs {
			fmt.Println(msgs)
			var msgBody models.ReciveMessage
			err := json.Unmarshal([]byte(msg.Body), &msgBody)
			if err != nil {
				fmt.Printf("Failed to parse message: %v", err)
				msg.Nack(false, true)
				continue
			}

			fmt.Printf("%+v\n", msgBody)
			jobType := msgBody.JobType

			switch jobType {
			case "upsert-testcase":
				service.AddAndUpdateTestCase(ch,db_pool,msg,msgBody,pubsub)
				msg.Ack(true)
			case "exercise-submit":
				service.RunSubmission(ch,db_pool,msg,msgBody,pubsub)
			default:
				fmt.Printf("job_type not wtf")
				msg.Ack(true)
			}
			msg.Ack(false)
		}
	}
}