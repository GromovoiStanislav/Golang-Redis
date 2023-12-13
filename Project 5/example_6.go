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

	//////////////////////////////////////////////////////////////////////////////////
	ctx := context.Background()
	
	// Выполнение команды FLUSHALL
	_ = client.FlushAll(ctx).Err()
	


	/*
		Incr
	*/


	client.Incr(ctx, "pipeline_counter")
	incr := client.Incr(ctx, "pipeline_counter")
	client.IncrBy(ctx, "pipeline_counter",5)
	client.Expire(ctx, "pipeline_counter", time.Hour)

	

	
	val, err := client.Get(ctx, "pipeline_counter").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("pipeline_counter", val) // pipeline_counter 7  String!!!
	


	fmt.Println(incr.Val()) //2 int!!!
	incr.SetVal(44)
	fmt.Println(incr.Val()) //44
	fmt.Println(incr.String()) //incr pipeline_counter: 44
	res,err := incr.Result() 
	if err != nil {
        panic(err)
    }
	fmt.Println(res) // 44 int!!!
}
