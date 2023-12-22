package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

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


	// Создаем WaitGroup для ожидания завершения горутин
	var wg sync.WaitGroup

	// Запускаем горутину для чтения сообщений из Redis и вывода их в консоль
	wg.Add(1)
	go readMessagesFromRedis(redisClient, &wg)

	// Ожидаем сигнала завершения (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Ожидаем сигнала завершения
	<-c

	// Завершаем горутину после получения сигнала
	wg.Done()

	// Закрываем соединение с Redis
	redisClient.Close()
	
}

func readMessagesFromRedis(client *redis.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// Блокирующе читаем сообщение из списка
		message, err := client.BRPop(context.Background(), 0, listKey).Result()
		if err != nil {
			log.Println("Ошибка при чтении сообщения из Redis:", err)
			return
		}

		// Выводим сообщение в консоль
		log.Printf("Принято сообщение: %v\n", message)
	}
}