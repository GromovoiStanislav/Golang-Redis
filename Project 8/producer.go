package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"redis-example/internal"
	"redis-example/internal/redis"
)


func main() {

	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	// Чтение переменной окружения REDIS_URL
	redisURL := os.Getenv("REDIS_URL")



	rdb:= redis.NewRedis(redisURL)
	defer rdb.Close()


	// Создаем экземпляр Task с использованием созданного клиента Redis
	taskRepo := redis.NewTask(rdb)


	// Пример создания задачи
	newTask := internal.Task{
		ID:        "14",
		Description:      "Sample Task",
		Priority: 4,
		Dates: internal.Dates{
			Start: time.Now(),
			Due: time.Now().Add(1 * time.Hour),
		},
		Categories: []internal.Category{"IT"},
		IsDone: false,
		// Дополнительные поля вашей задачи...
	}
	newTask.Validate()
	if err != nil {
		log.Println(err)
	} else {
		
		err := taskRepo.Created(context.Background(), newTask)
		if err != nil {
			log.Fatalf("Failed to publish task created event: %v", err)
		}
		
	}


	time.Sleep(3*time.Second)

	// Пример удаления задачи
	taskIDToDelete := "2"
	err = taskRepo.Deleted(context.Background(), taskIDToDelete)
	if err != nil {
		log.Fatalf("Failed to publish task deleted event: %v", err)
	}

	log.Println("Task events published successfully!")

}
