package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"redis-example/internal/utility"
)

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

	// Set "FOO" to be associated with "BAR"
	rdb.Set(ctx, "FOO", "BAR", 0)

	result, err := rdb.Get(ctx, "FOO").Result()
	if err != nil {
		fmt.Println("FOO not found")
	} else {
		fmt.Printf("FOO has value %s\n", result)
	}

	// Deleting the key "FOO" and its associated value
	rdb.Del(ctx, "FOO")

	result, err = rdb.Get(ctx, "FOO").Result()
	if err != nil {
		fmt.Println("FOO not found")
	} else {
		fmt.Printf("FOO has value %s\n", result)
	}
}