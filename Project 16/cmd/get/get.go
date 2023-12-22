package main

import (
	"context"
	"fmt"

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

	result, err := rdb.Get(ctx, "FOO").Result()
	if err != nil {
		fmt.Println("Key FOO not found in Redis cache")
	} else {
		fmt.Printf("FOO has value %s\n", result)
	}

	intValue, err := rdb.Get(ctx, "INT").Int()
	if err != nil {
		fmt.Println("Key INT not found in Redis cache")
	} else {
		fmt.Printf("INT has value %d\n", intValue)
	}

	var person Person
	err = rdb.HGetAll(ctx, "STRUCT").Scan(&person)
	if err != nil {
		fmt.Println("Key STRUCT not found in Redis cache")
	} else {
		fmt.Printf("STRUCT has value %+v\n", person)
	}

	result, err = rdb.Get(ctx, "BAZ").Result()
	if err != nil {
		fmt.Println("Key BAZ not found in Redis cache")
	} else {
		fmt.Printf("BAZ has value %s\n", result)
	}
}