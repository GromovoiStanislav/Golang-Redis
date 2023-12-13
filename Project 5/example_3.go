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
		list
	*/
	
	// Заполняем список значениями
	listKey := "list"
	values := []int{3, 1, 2,6,5} // Пример значений, замените на ваши

	// RPUSH - добавление значений в конец списка
	for _, value := range values {
		if err := client.RPush(ctx, listKey, value).Err(); err != nil {
			log.Fatal(err)
		}
	}

	// Опционально: выводим значения списка
	listValues, err := client.LRange(ctx, listKey, 0, -1).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Список после заполнения:", listValues) //: [3 1 2 6 5]

	// Теперь можем использовать SORT для сортировки и извлечения элементов
	vals, err := client.Sort(ctx, listKey, &redis.Sort{
		Offset: 0,
		Count:  2,
		Order:  "ASC",
	}).Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Отсортированные значения:", vals)//: [1 2]
	
}
