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

	// Perform basic diagnostic to check if the connection is working
	// Expected result > ping: PONG
	// If Redis is not running, error case is taken instead
	status, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Redis connection was refused")
		return
	}
	fmt.Println(status)
}