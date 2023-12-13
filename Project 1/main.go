package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	// Создаем клиент Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Адрес Redis-сервера
		Password: "",               // Пароль (если требуется)
		DB:       0,                // Номер базы данных
	})

	// Подписываемся на канал
	pubsub := client.Subscribe(ctx, "mychannel")
	defer pubsub.Close()

	// Канал для получения сообщений
	ch := pubsub.Channel()

	// Горутина для обработки сообщений
	go func() {
		for msg := range ch {
			fmt.Printf("Получено сообщение из канала: %s\n", msg.Payload)
		}
	}()

	// Отправляем сообщения в канал
	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Сообщение %d", i)
		err := client.Publish(ctx, "mychannel", message).Err()
		if err != nil {
			log.Printf("Ошибка отправки сообщения: %v\n", err)
		} else {
			fmt.Printf("Отправлено сообщение: %s\n", message)
		}

		// Дадим немного времени для обработки сообщения
		time.Sleep(time.Second)
	}

	// Дожидаемся завершения работы
	//select {}
}

// go run main.go
