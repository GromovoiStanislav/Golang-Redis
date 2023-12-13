package main

import (
	"context"
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

	//////////////////////////////////////////////////////////////////////////////////
	ctx := context.Background()
	
	// Выполнение команды FLUSHALL
	_ = client.FlushAll(ctx).Err()

	/*
	Delete keys without a ttl
	This example demonstrates how to use SCAN and pipelines to efficiently delete keys without a TTL
	*/

	iter := client.Scan(ctx, 0, "", 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()

		d, err := client.TTL(ctx, key).Result()
		if err != nil {
			panic(err)
		}

		if d == -1 { // -1 means no TTL
			if err := client.Del(ctx, key).Err(); err != nil {
				panic(err)
			}
		}
	}

	if err := iter.Err(); err != nil {
		panic(err)
	}
}
