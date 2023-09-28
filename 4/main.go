package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
	"strconv"
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

	fmt.Println("=================== Example ===========================")

	/// User #1
	type User struct {
		ID   string `redis:"id"`
		Name string `redis:"name"`
		Age  int    `redis:"age"`
	}

	user := User{
		ID:   "001",
		Name: "Tom",
		Age:  30,
	}

	// Преобразуем структуру User в map
	userFields := map[string]interface{}{
		"id":   user.ID,
		"name": user.Name,
		"age":  user.Age,
	}

	// Используем мультикоманду для выполнения нескольких операций атомарно
	pipe := client.Pipeline()
	pipe.HMSet(ctx, "users:"+user.ID, userFields)
	pipe.RPush(ctx, "users", user.ID)
	// Выполняем мультикоманду
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Ошибка при выполнении мультикоманды: %v", err)
	}

	/// User #2
	user = User{
		ID:   "002",
		Name: "Toma",
		Age:  30,
	}

	// Преобразуем структуру User в map
	userFields = map[string]interface{}{
		"id":   user.ID,
		"name": user.Name,
		"age":  user.Age,
	}

	// Используем мультикоманду для выполнения нескольких операций атомарно
	pipe = client.Pipeline()
	pipe.HMSet(ctx, "users:"+user.ID, userFields)
	pipe.RPush(ctx, "users", user.ID)
	// Выполняем мультикоманду
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Ошибка при выполнении мультикоманды: %v", err)
	}

	/// User #3
	user = User{
		ID:   "003",
		Name: "Tomas",
		Age:  30,
	}

	// Преобразуем структуру User в map
	userFields = map[string]interface{}{
		"id":   user.ID,
		"name": user.Name,
		"age":  user.Age,
	}

	// Используем мультикоманду для выполнения нескольких операций атомарно
	pipe = client.Pipeline()
	pipe.HMSet(ctx, "users:"+user.ID, userFields)
	pipe.RPush(ctx, "users", user.ID)
	// Выполняем мультикоманду
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Fatalf("Ошибка при выполнении мультикоманды: %v", err)
	}

	///////////////////////////  READ All to Struct  //////////////////////

	// Получаем список пользователей из списка "users"
	userIDs, err := client.LRange(ctx, "users", 0, -1).Result()
	if err != nil {
		log.Fatalf("Ошибка LRange: %v", err)
	}

	// Создаем слайс для хранения результатов
	var users []User

	// Используем цикл для получения данных о каждом пользователе
	for _, userID := range userIDs {
		userFields, err := client.HGetAll(ctx, "users:"+userID).Result()
		if err != nil {
			log.Printf("Ошибка HGetAll: %v", err)
			continue // Пропустить ошибочные записи и продолжить цикл
		}

		var user User

		// Преобразуем map в структуру User вручную
		user.ID = userFields["id"]
		user.Name = userFields["name"]
		ageStr := userFields["age"]
		user.Age, _ = strconv.Atoi(ageStr)

		users = append(users, user)
	}

	fmt.Println("Содержимое списка пользователей:")
	for _, user := range users {
		fmt.Println("User:", user)
	}

	//////////////////// READ one user ///////////////////////

	id := "001"
	userID := fmt.Sprintf("users:%s", id)

	// Получаем индекс элемента в списке "users"
	_, err = client.LPos(ctx, "users", id, redis.LPosArgs{}).Result()
	if err != nil {
		log.Printf("Ошибка LPos: %v", err)
		fmt.Printf("Пользователь с индексом %s не найден\n", id)
	} else {

		userFields, err := client.HGetAll(ctx, userID).Result()
		if err != nil {
			log.Printf("Ошибка HGetAll: %v", err)
		}

		// Выводим данные о пользователе
		fmt.Println("Данные о пользователе:")
		fmt.Println(userFields)
	}

	////////////////////// UPDATE one user ///////////////////////

	id = "002"
	userID = fmt.Sprintf("users:%s", id)

	// Получаем индекс элемента в списке "users"
	_, err = client.LPos(ctx, "users", id, redis.LPosArgs{}).Result()
	if err != nil {
		log.Printf("Ошибка LPos: %v", err)
		fmt.Printf("Пользователь с индексом %s не найден\n", id)
	} else {

		// Создаем новый пайплайн
		pipe = client.Pipeline()
		pipe.HSet(ctx, userID, "city", "Moscow")
		pipe.HSet(ctx, userID, "age", 31)

		// Выполняем пайплайн
		_, err = pipe.Exec(ctx)
		if err != nil {
			log.Printf("Ошибка при выполнении пайплайна: %v", err)
		}

	}

	////////////////////// Delete one user ///////////////////////

	id = "003"
	userID = fmt.Sprintf("users:%s", id)

	// Получаем индекс элемента в списке "users"
	_, err = client.LPos(ctx, "users", id, redis.LPosArgs{}).Result()
	if err != nil {
		log.Printf("Ошибка LPos: %v", err)
		fmt.Printf("Пользователь с индексом %s не найден\n", id)
	} else {

		// Создаем новый пайплайн
		pipe = client.Pipeline()
		pipe.LRem(ctx, "users", 0, id)
		pipe.Del(ctx, userID)

		// Выполняем пайплайн
		_, err = pipe.Exec(ctx)
		if err != nil {
			log.Printf("Ошибка при выполнении пайплайна: %v", err)
		}

	}

	///////////////////////////  READ All to Map  //////////////////////

	// Получаем список пользователей из списка "users"
	userIDs, err = client.LRange(ctx, "users", 0, -1).Result()
	if err != nil {
		log.Fatalf("Ошибка LRange: %v", err)
	}

	// Создаем слайс для хранения результатов
	var users_ []map[string]string

	// Используем цикл для получения данных о каждом пользователе
	for _, userID := range userIDs {
		userFields, err := client.HGetAll(ctx, "users:"+userID).Result()
		if err != nil {
			log.Printf("Ошибка HGetAll: %v", err)
			continue // Пропустить ошибочные записи и продолжить цикл
		}

		users_ = append(users_, userFields)
	}

	fmt.Println("Содержимое списка пользователей:")
	for _, user := range users_ {
		fmt.Println("User:", user)
	}

	fmt.Println("=================== Exit ===========================")
}
