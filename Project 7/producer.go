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

	
	// Публикация сообщений в стрим "mystream"
	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Message %d", i)

		// Добавление сообщения в стрим "mystream"
		_, err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: "mystream",
			Values: map[string]interface{}{"message": message},
		}).Result()

		if err != nil {
			fmt.Println("Error publishing message:", err)
			return
		}

		time.Sleep(time.Second) // Задержка между сообщениями
	}

	fmt.Println("Producer finished")


}
