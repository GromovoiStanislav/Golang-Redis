package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"redis-example/internal/utility"
)

type Person struct {
	Name string `redis:"name"`
	Age  int    `redis:"age"`
}

func main() {
	ctx := context.Background()
	// Ensure that you have Redis running on your system
	rdb := redis.NewClient(&redis.Options{
		Addr:     utility.Address(),
		Username: utility.Username(),
		Password: utility.Password(),
		DB:       utility.Database(),
	})
	// Ensure that the connection is properly closed gracefully
	defer rdb.Close()

	_, err := rdb.Set(ctx, "FOO", "BAR", 0).Result()
	if err != nil {
		fmt.Println("Failed to add FOO <> BAR key-value pair")
		return
	}
	rdb.Set(ctx, "INT", 5, 0)
	rdb.Set(ctx, "FLOAT", 5.5, 0)
	rdb.Set(ctx, "EXPIRING", 15, 30*time.Minute)
	rdb.Set(ctx, "LIST", []string{"Hello"}, 0)

	rdb.HSet(ctx, "STRUCT", Person{"John Doe", 15})
}