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


	var client *redis.Client
	if redisURL == "" {
		client = redis.NewClient(&redis.Options{
			Addr:	  "localhost:6379",
			Password: "", // no password set
			DB:		  0,  // use default DB
		})
	} else {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(err)
		}
		client = redis.NewClient(opt)
	}
	defer client.Close()


	ctx := context.Background()

	
	// Создание группы для стрима "mystream"
	//client.XGroupDestroy(ctx, "mystream", "group-2")
	client.XGroupCreate(ctx, "mystream", "group-2", "0")
	client.XGroupCreateConsumer(ctx, "mystream", "group-2", "consumer-1")


	// Чтение сообщений из стрима "mystream" группой "mygroup"
	for {
		// Чтение сообщений из стрима "mystream" группой "mygroup"
		streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "group-2",
			Consumer: "consumer-1",
			Streams:  []string{"mystream", ">"},
			Block:    0,
			Count:    1,
			NoAck:    true,
		}).Result()

		if err != nil {
			fmt.Println("Error reading from stream:", err)
			return
		}

		// Обработка полученных сообщений
		for _, stream := range streams {
			for _, message := range stream.Messages {
				fmt.Printf("Stream: %s, ID: %s, Message: %s\n", stream.Stream, message.ID, message.Values["message"])

				// Подтверждение получения сообщения
				ackCmd := client.XAck(ctx, "mystream", "group-2", message.ID)
				if err := ackCmd.Err(); err != nil {
					fmt.Println("Error acknowledging message:", err)
					return
				}
				fmt.Printf("Message %s acknowledged\n", message.ID)
			}
		}

		time.Sleep(time.Second) // Задержка между попытками чтения
	}


}
