package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
)

var ctx = context.Background()

func main() {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	// Чтение переменной окружения REDIS_URL
	redisURL := os.Getenv("REDIS_URL")
	// Если REDIS_URL не установлена, используйте значение по умолчанию "localhost:6379"
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	// Разбор URL-адреса Redis
	parsedURL, err := url.Parse(redisURL)
	if err != nil {
		log.Fatalf("Ошибка разбора URL-адреса Redis: %v", err)
	}

	// Извлечение компонент URL-адреса
	hostname := parsedURL.Hostname()
	port := parsedURL.Port()
	username := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()

	// Если хост и порт не указаны, используйте значения по умолчанию
	if hostname == "" {
		hostname = "localhost"
	}
	if port == "" {
		port = "6379"
	}

	// Создаем клиент Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", hostname, port), // Хост и порт
		Password: password,                             // Пароль (если требуется)
		Username: username,                             // username (если требуется)
		DB:       0,                                    // Номер базы данных
	})

	defer client.Close()

	//////////////////////////////////////////////////////////////////////////////////

	// Выполнение команды FLUSHALL
	err = client.FlushAll(ctx).Err()
	if err != nil {
		log.Printf("Ошибка выполнения команды FLUSHALL: %v\n", err)
	} else {
		fmt.Println("Команда FLUSHALL выполнена успешно.")
	}

	// Определение структуры bike1
	type Bike struct {
		Model       string
		Brand       string
		Price       int
		Type        string
		Description string
		Specs       struct {
			Material string
			Weight   float64
		}
	}

	bike1 := Bike{
		Model:       "Hyperion",
		Brand:       "Velorim",
		Price:       844,
		Type:        "Enduro bikes",
		Description: "This is a mid-travel trail slayer...",
		Specs: struct {
			Material string
			Weight   float64
		}{
			Material: "full-carbon",
			Weight:   8.7,
		},
	}

	log.Println("=================== JSON->String ======================")

	// Преобразуем структуру в JSON
	bike1JSON, err := json.Marshal(bike1)
	if err != nil {
		log.Fatalf("Ошибка преобразования в JSON: %v", err)
	}

	// Сохраняем bike1 в Redis
	err = client.Set(ctx, "bikes:1", bike1JSON, 0).Err()
	if err != nil {
		log.Fatalf("Ошибка Set: %v", err)
	}

	// Получаем данные и декодируем их из JSON
	bike1Bytes, err := client.Get(ctx, "bikes:1").Bytes()
	if err != nil {
		log.Fatalf("Ошибка Get: %v", err)
	}

	var retrievedBike Bike
	err = json.Unmarshal(bike1Bytes, &retrievedBike)
	if err != nil {
		log.Fatalf("Ошибка декодирования из JSON: %v", err)
	}

	fmt.Printf("Модель: %s, Материал: %s\n", retrievedBike.Model, retrievedBike.Specs.Material)

	log.Println("=================== JSON->HASH ======================")

	// Преобразуем структуру в хеш и сохраняем все поля
	bikeFields := map[string]interface{}{
		"model":          bike1.Model,
		"brand":          bike1.Brand,
		"price":          bike1.Price,
		"type":           bike1.Type,
		"description":    bike1.Description,
		"specs.material": bike1.Specs.Material,
		"specs.weight":   bike1.Specs.Weight,
	}

	err = client.HMSet(ctx, "bikes:2", bikeFields).Err()
	if err != nil {
		log.Fatalf("Ошибка HMSet: %v", err)
	}

	// Получаем все поля и значения хеша с помощью HGETALL
	bikeData, err := client.HGetAll(ctx, "bikes:2").Result()
	if err != nil {
		log.Fatalf("Ошибка HGetAll: %v", err)
	}

	// Выводим все поля и их значения
	for field, value := range bikeData {
		fmt.Printf("%s: %s\n", field, value)
	}

	// Изменяем некоторые поля
	err = client.HSet(ctx, "bikes:2", "model", "Hyperion2").Err()
	if err != nil {
		log.Fatalf("Ошибка HSet: %v", err)
	}

	err = client.HSet(ctx, "bikes:2", "brand", "GIGANT").Err()
	if err != nil {
		log.Fatalf("Ошибка HSet: %v", err)
	}

	// ... Измените остальные поля аналогичным образом

	// Получаем данные с помощью HGET
	model, err := client.HGet(ctx, "bikes:2", "model").Result()
	if err != nil {
		log.Fatalf("Ошибка HGet: %v", err)
	}

	brand, err := client.HGet(ctx, "bikes:2", "brand").Result()
	if err != nil {
		log.Fatalf("Ошибка HGet: %v", err)
	}

	// ... Получите остальные поля аналогичным образом

	// Выводим полученные данные
	fmt.Printf("Модель: %s\n", model)
	fmt.Printf("Бренд: %s\n", brand)

	fmt.Println("===================Exit===========================")
}
