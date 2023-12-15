package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)



func main() {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	// Чтение переменной окружения REDIS_URL
	redisURL := os.Getenv("REDIS_URL")


	var rdb *redis.Client
	if redisURL == "" {
		rdb = redis.NewClient(&redis.Options{
			Addr:	  "localhost:6379",
			Password: "", // no password set
			DB:		  0,  // use default DB
		})
	} else {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(err)
		}
		rdb = redis.NewClient(opt)
	}
	defer rdb.Close()


	ctx := context.Background()
	


	// To subscribe to a channel and receive a message:
	go func() {
		pubsub := rdb.Subscribe(ctx, "mychannel-1")
		defer pubsub.Close()
		channel := pubsub.Channel()

		defer func() {
			pubsub.Close()
		}()

		// Listen for messages
		for msg:= range channel {
			fmt.Println("Subscribe",msg.Channel, msg.Payload)
		}
	}()

	go func() {
		pubsub := rdb.PSubscribe(ctx, "mychannel*")
		defer pubsub.Close()
		channel := pubsub.Channel()

		defer func() {
			pubsub.Close()
		}()

		// Listen for messages
		for msg:= range channel {
			fmt.Println("PSubscribe",msg.Channel, msg.Payload)
		}
	}()



	
	// Запускаем бесконечный цикл, который будет ждать событий от таймера
	ticker := time.NewTicker(1 * time.Second)
	for {
		// Ожидаем события от таймера
		<-ticker.C

		rdb.Publish(ctx, "mychannel-1", "payload")
	}
		
}
