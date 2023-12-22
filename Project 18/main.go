package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)


var redisClient *redis.Client

var ctx = context.Background()


func main() {
	// Connecting to a Redis Instance
	redisClient = createRedisClient()
	defer redisClient.Close()

	// Проверка подключения к Redis
	err := redisClient.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to Redis:", err)
		return
	}
	log.Println("Redis is working....")

	hllandset()
}


// Connecting to a Redis Instance
func createRedisClient() *redis.Client { 

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	redisURL := os.Getenv("REDIS_URL")

	if redisURL == "" {
		redisClient := redis.NewClient(&redis.Options{
			Addr:	  "localhost:6379",
			Password: "", // no password set
			DB:		  0,  // use default DB
		})
		return redisClient

	} else {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(err)
		}
		redisClient := redis.NewClient(opt)
		return redisClient
	}
} 

func hllandset() {
	rand.Seed(time.Now().UnixNano())

	// Создание объекта Pipeline
	pipe := redisClient.Pipeline()

	// Генерация случайных чисел и добавление их в множество и HyperLogLog
	numbers := []string{}
	for i := 0; i < 100_000; i++ {
		num := rand.Intn(100) // генерация случайного числа от 0 до 99
		numbers = append(numbers, strconv.Itoa(num))
	}

	pipe.SAdd(ctx, "set_of_numbers", numbers)
	pipe.PFAdd(ctx, "hll_of_numbers", numbers)

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Println("Error executing pipeline:", err)
		return
	}


	// Вывод количества уникальных элементов в множестве и HyperLogLog
	log.Println("Number of unique elements (SCARD):", redisClient.SCard(ctx, "set_of_numbers").Val())
	log.Println("Number of unique elements (PFCOUNT):", redisClient.PFCount(ctx, "hll_of_numbers").Val())



	// // Объединение двух HyperLogLog в один с использованием PFMerge.
	// result := redisClient.PFMerge(ctx, destinationKey, key1, key2)
	// if result.Err() != nil {
	// 	log.Println("Error merging HyperLogLogs:", result.Err())
	// 	return
	// }
}