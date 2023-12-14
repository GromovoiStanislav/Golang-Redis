package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	// // Чтение всех сообщений из стрима "mystream"
	// streams, err := client.XRead(ctx, &redis.XReadArgs{
	// 	Streams: []string{"mystream", "0"},
	// 	Count:   0, // Количество сообщений для чтения (указывается на усмотрение)
	// 	Block:   0,
	// }).Result()
	streams, err := client.XReadStreams(ctx,"mystream","$").Result()
	if err != nil {
		fmt.Println("Error reading from stream:", err)
		return
	}
	// Обработка полученных сообщений
	for _, stream := range streams {
		for _, message := range stream.Messages {
			fmt.Printf("Received message: %s\n", message.Values["message"])
		}
	}


	// Получение всех сообщений из стрима "mystream"
	//messages, err := client.XRange(ctx, "mystream", "-", "+").Result()
	//messages, err := client.XRange(ctx, "mystream", "1702541071069-0", "1702541152169-0").Result()
	//messages, err := client.XRangeN(ctx, "mystream", "-", "+",2).Result()
	//messages, err := client.XRangeN(ctx, "mystream", "1702541071069-0", "1702541152169-0", 2).Result()
	//messages, err := client.XRevRange(ctx, "mystream", "+", "-").Result()
	//messages, err := client.XRevRange(ctx, "mystream", "1702541152169-0", "1702541071069-0").Result()
	//messages, err := client.XRevRangeN(ctx, "mystream", "+", "-",3).Result()
	// messages, err := client.XRevRangeN(ctx, "mystream", "1702541152169-0", "1702541071069-0", 3).Result()
	// if err != nil {
	// 	fmt.Println("Error getting messages:", err)
	// 	return
	// }
	// // Обработка полученных сообщений
	// for _, message := range messages {
	// 	fmt.Printf("ID: %s, Message: %s\n", message.ID, message.Values["message"])
	// }


}
