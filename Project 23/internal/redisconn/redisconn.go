package redisconn

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)


var redisOnce sync.Once
var redisInstance *redis.Client

func GetRedisConnection() *redis.Client {
	
	redisOnce.Do(func() {
		if err := godotenv.Load(); err != nil {
			panic("Ошибка загрузки файла .env: " + err.Error())
		}
		redisURL := os.Getenv("REDIS_URL")
	
		if redisURL == "" {
			redisInstance = redis.NewClient(&redis.Options{
				Addr:	  "localhost:6379",
				Password: "",
				DB:		  0,
			})
		} else {
			opt, err := redis.ParseURL(redisURL)
			if err != nil {
				panic(err)
			}
			redisInstance= redis.NewClient(opt)
		}
	})

    return redisInstance
}