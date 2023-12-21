package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

var redisClient *redis.Client
var ctx = context.Background()

func main() {
	connectRedis(ctx)
	defer redisClient.Close()

    subscriber := redisClient.Subscribe(ctx, "send-user-data")

    user := User{}

    for {
        msg, err := subscriber.ReceiveMessage(ctx)
        if err != nil {
            panic(err)
        }

        if err := json.Unmarshal([]byte(msg.Payload), &user); err != nil {
            panic(err)
        }

        fmt.Println("Received message from " + msg.Channel + " channel.")
        fmt.Printf("%+v\n", user)
    }
}

func connectRedis(ctx context.Context) {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	redisURL := os.Getenv("REDIS_URL")


	if redisURL == "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:	  "localhost:6379",
			Password: "", // no password set
			DB:		  0,  // use default DB
		})
	} else {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(err)
		}
		redisClient = redis.NewClient(opt)
	}
}