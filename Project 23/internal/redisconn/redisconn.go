package redisconn

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func GetRedisConnection() *redis.Client {

    if err := godotenv.Load(); err != nil {
		panic("Ошибка загрузки файла .env: " + err.Error())
	}
	redisURL := os.Getenv("REDIS_URL")

	log.Println("GetRedisConnection")

	if redisURL == "" {
		return redis.NewClient(&redis.Options{
			Addr:	  "localhost:6379",
			Password: "",
			DB:		  0,
		})
	} else {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(err)
		}
		return redis.NewClient(opt)
	}
}