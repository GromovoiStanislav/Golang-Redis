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
	 Store and retrieve a map
	*/
	{
		session := map[string]string{"name": "John", "surname": "Smith", "company": "Redis", "age": "29"}
		for k, v := range session {
			err := client.HSet(ctx, "user-session:123", k, v).Err()
			if err != nil {
				panic(err)
			}
		}

		userSession := client.HGetAll(ctx, "user-session:123").Val()
		fmt.Println(userSession) // map[age:29 company:Redis name:John surname:Smith]
	}


	{
		// Пример данных пользователя
		userID := 125
		userData := map[string]interface{}{
			"username": "john_doe",
			"email":    "john@example.com",
			"company": "Redis",
			 "age": "30",
			// Другие данные пользователя...
		}

		// Добавляем данные пользователя в хэш-множество
		err := client.HMSet(ctx, fmt.Sprintf("user-session:%d", userID), userData).Err()
		if err != nil {
			fmt.Println("Ошибка при добавлении данных пользователя:", err)
			return
		}

		userSession := client.HGetAll(ctx, "user-session:125").Val()
		fmt.Println(userSession) // map[age:30 company:Redis email:john@example.com username:john_doe]
	}

	// Example for scanning hash fields into a struct
	{
		type Model struct {
			Str1    string   `redis:"str1"`
			Str2    string   `redis:"str2"`
			Int     int      `redis:"int"`
			Bool    bool     `redis:"bool"`
			Ignored struct{} `redis:"-"`
		}

		if _, err := client.Pipelined(ctx, func(rdb redis.Pipeliner) error {
			rdb.HSet(ctx, "key", "str1", "hello")
			rdb.HSet(ctx, "key", "str2", "world")
			rdb.HSet(ctx, "key", "int", 123)
			rdb.HSet(ctx, "key", "bool", 1)
			return nil
		}); err != nil {
			panic(err)
		}

		var model1 Model
		// Scan all fields into the model.
		if err := client.HGetAll(ctx, "key").Scan(&model1); err != nil {
			panic(err)
		}
		fmt.Println(model1) // {hello world 123 true {}}

		var model2 Model
		// Scan a subset of the fields.
		if err := client.HMGet(ctx, "key", "str1", "int").Scan(&model2); err != nil {
			panic(err)
		}
		fmt.Println(model2) // {hello world 123 true {}}
	}
	

	// Iterating over keys
	{
		iter := client.Scan(ctx, 0, "user-session:*", 0).Iterator()
		for iter.Next(ctx) {
			fmt.Println("keys", iter.Val())
		}
		//keys user-session:125
		//keys user-session:123
		if err := iter.Err(); err != nil {
			panic(err)
		}
	}


	{
		// Добавляем данные в множество "hash-key"
		client.HSet(ctx, "hash-key", "user:session:55", 55)
		client.HSet(ctx, "hash-key", "user:session:56", 56)
		client.HSet(ctx, "hash-key", "user:session:57", 57)
		client.HSet(ctx, "hash-key", "user:data:58", 58)
	
		data := client.HGetAll(ctx, "hash-key").Val()
		fmt.Println(data) //map[user:data:58:1 user:session:55:1 user:session:56:1 user:session:57:1]

		iter := client.HScan(ctx, "hash-key", 0, "user:session:*", 0).Iterator()
		for iter.Next(ctx) {
			fmt.Println("hash-key", iter.Val())
		}
		if err := iter.Err(); err != nil {
			panic(err)
		}
	}


	{
		// Set hash field-values
		err = client.HSet(ctx, "user:1", map[string]interface{}{
			"name":  "John Doe",
			"email": "john@example.com",
			"age":   25,
		}).Err()
		if err != nil {
			log.Fatal(err)
		}
		
		// Get hash field-values
		userInfo, err := client.HGetAll(ctx, "user:1").Result()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("User Info:", userInfo)
	}
}
