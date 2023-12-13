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

	//////////////////////////////////////////////////////////////////////////////////
	ctx := context.Background()
	
	/*
	Redis PubSub
	*/

	// To subscribe to a channel and receive a message:
	go func() {
		pubsub := rdb.Subscribe(ctx, "mychannel1")
		defer pubsub.Close()
		
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				panic(err)
			}
		
			fmt.Println(msg.Channel, msg.Payload)
		}
	}()


	go func() {
		pubsub := rdb.Subscribe(ctx, "mychannel2")
		defer pubsub.Close()
		
		ch := pubsub.Channel()
		for msg := range ch {
			fmt.Println(msg.Channel, msg.Payload)
		}
	}()


	

	// Создаем новый таймер, который срабатывает каждую секунду
	ticker := time.NewTicker(1 * time.Second)

	// Запускаем бесконечный цикл, который будет ждать событий от таймера
	for {
		// Ожидаем события от таймера
		<-ticker.C

		// To publish a message:
		err = rdb.Publish(ctx, "mychannel1", time.Now()).Err()
		if err != nil {
			panic(err)
		}

		err = rdb.Publish(ctx, "mychannel2", "payload").Err()
		if err != nil {
			panic(err)
		}
	}
		
}
