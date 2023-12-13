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

	//////////////////////////////////////////////////////////////////////////////////
	ctx := context.Background()
	
	// Выполнение команды FLUSHALL
	_ = client.FlushAll(ctx).Err()
	

	/*
		HyperLogLog
	*/

	for i := 0; i < 10; i++ {
		if err := client.PFAdd(ctx, "myset", fmt.Sprint(i)).Err(); err != nil {
			panic(err)
		}
	}

	card, err := client.PFCount(ctx, "myset").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("set cardinality", card) //set cardinality 10
	



}
