package main

import (
	"context"
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
		Pipeline
	*/

	{
		// open the pipeline
		pipe := client.Pipeline()

		// submit some commands

		client.Set(ctx, "test_key", "test_value", 0)
		client.Set(ctx, "food", "cheese", redis.KeepTTL)

		pipe.Incr(ctx, "pipeline_counter")
		incr := pipe.Incr(ctx, "pipeline_counter")
		pipe.IncrBy(ctx, "pipeline_counter",5)
		pipe.Expire(ctx, "pipeline_counter", time.Hour)

		// execute with trace
		_, err := pipe.Exec(ctx)
		if err != nil {
			panic(err)
		}


		// The value is available only after Exec is called.
		val, err := client.Get(ctx, "test_key").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("test_key", val) // test_key test_value

		val, err = client.Get(ctx, "food").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("food", val) // food cheese

		val, err = client.Get(ctx, "pipeline_counter").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("pipeline_counter", val) // pipeline_counter 7  String!!!

		fmt.Println(incr.Val()) // 2 int!!!
	}
	

	// Alternatively, you can use Pipelined which calls Exec when the function exits:
	{
		var incr *redis.IntCmd

		_, err := client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			
			client.Set(ctx, "test_key2", "test_value", 0)
			client.Set(ctx, "food2", "cheese", redis.KeepTTL)

			pipe.Incr(ctx, "pipeline_counter2")
			incr = pipe.Incr(ctx, "pipeline_counter2")
			pipe.IncrBy(ctx, "pipeline_counter2",5)
			pipe.Expire(ctx, "pipeline_counter2", time.Hour)

			pipe.Expire(ctx, "pipelined_counter2", time.Hour)
			
			return nil
		})
		if err != nil {
			panic(err)
		}
		
		// The value is available only after the pipeline is executed.
		val, err := client.Get(ctx, "test_key2").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("test_key2", val) // test_key test_value

		val, err = client.Get(ctx, "food2").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("food2", val) // food cheese

		val, err = client.Get(ctx, "pipeline_counter2").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("pipeline_counter", val) // pipeline_counter 7  String!!!

		fmt.Println(incr.Val())
	}


	// Pipelines also return the executed commands so can iterate over them to retrieve results:
	{
		cmds, err := client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			// for i := 0; i < 100; i++ {
			// 	pipe.Get(ctx, fmt.Sprintf("key%d", i))
			// }

			pipe.Get(ctx, "test_key")
			pipe.Get(ctx, "test_key2")
			pipe.Get(ctx, "food")
			pipe.Get(ctx, "food2")
			pipe.Get(ctx, "pipeline_counter")
			pipe.Get(ctx, "pipeline_counter2")
			
			return nil
		})
		if err != nil {
			panic(err)
		}
		
		for _, cmd := range cmds {
			fmt.Println(cmd.(*redis.StringCmd).Val())
		}
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
}
