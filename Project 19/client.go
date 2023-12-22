package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)


const listKey = "messageList"

func main() {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	// Чтение переменной окружения REDIS_URL
	redisURL := os.Getenv("REDIS_URL")


	var redisClient *redis.Client
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
	defer redisClient.Close()

	// Generate a new background context that  we will use
	ctx := context.Background()
	

	// Проверяем подключение к Redis
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to Redis:", err)
		return
	}
	log.Println("Redis is working....")

	// Открываем стандартный ввод для чтения сообщений из консоли
	scanner := bufio.NewScanner(os.Stdin)

	for {
		log.Print("Введите сообщение (или 'exit' для завершения): ")
		scanner.Scan()
		message := scanner.Text()

		if message == "exit" {
			break
		}

		// Отправляем сообщение в Redis
		err := pushMessageToRedis(redisClient, message)
		if err != nil {
			log.Println("Ошибка при отправке сообщения в Redis:", err)
		}
	}



}

func pushMessageToRedis(client *redis.Client, message string) error {
	// Добавляем сообщение в список с помощью LPUSH
	err := client.LPush(context.Background(), listKey, message).Err()
	return err
}