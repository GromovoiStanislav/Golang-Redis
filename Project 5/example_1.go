package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

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
		Store and retrieve a simple string
	*/

	err = client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val) // foo bar

	err = client.Set(ctx, "key", "value", 10*time.Second).Err()
    if err != nil {
        panic(err)
    }

    val, err = client.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key", val) // key value

    val2, err := client.Get(ctx, "key2").Result()
    if err == redis.Nil {
        fmt.Println("key2 does not exist") // key2 does not exist
    } else if err != nil {
        panic(err)
    } else {
        fmt.Println("key2", val2)
    }


	val, err = client.Get(ctx, "key2").Result()
	switch {
	case err == redis.Nil:
		fmt.Println("key2 does not exist")
	case err != nil:
		fmt.Println("Get failed", err)
	case val == "":
		fmt.Println("value is empty")
	}


	// SET key value EX 10 NX
	set, err := client.SetNX(ctx, "key1", "value", 10*time.Second).Result()
	fmt.Println("set", set) // set true

	// SET key value keepttl NX
	set, err = client.SetNX(ctx, "key1", "value", redis.KeepTTL).Result()
	fmt.Println("set", set) // set false

	// custom command
	res, err := client.Do(ctx, "set", "key3", "value").Result()
	fmt.Println("res", res) // res OK



	{
		type Author struct {
			Name string `json:"name"`
			Age int `json:"age"`
		}

		json, err := json.Marshal(Author{Name: "Elliot", Age: 25})
		if err != nil {
			fmt.Println(err)
		}

		err = client.Set(ctx, "id1234", json, 0).Err()
		if err != nil {
			fmt.Println(err)
		}
		val, err := client.Get(ctx, "id1234").Result()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(val)
	}


	// Iterating over keys
	{
		iter := client.Scan(ctx, 0, "*", 0).Iterator()
		for iter.Next(ctx) {
			fmt.Println("keys:", iter.Val())//keys user-session:125 //keys user-session:123
		}
		// keys: key1
		// keys: key
		// keys: foo
		// keys: key3
		if err := iter.Err(); err != nil {
			panic(err)
		}
	}

	{
		var cursor uint64
		for {
			var keys []string
			var err error
			keys, cursor, err = client.Scan(ctx, cursor, "*", 0).Result()
			if err != nil {
				panic(err)
			}

			for _, key := range keys {
				fmt.Println("key:", key)
			}
			// key: key1
			// key: key
			// key: foo
			// key: key3

			if cursor == 0 { // no more keys
				break
			}
		}
	}
}
