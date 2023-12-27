package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	client = &redisClient{}
	ctx = context.Background()
)

type redisClient struct {
	c *redis.Client
}

//GetClient get the redis client
func initialize() *redisClient {

	if err := godotenv.Load(); err != nil {
		panic("Ошибка загрузки файла .env: " + err.Error())
	}

	redisURL := os.Getenv("REDIS_URL")


	var c *redis.Client
	if redisURL == "" {
		c = redis.NewClient(&redis.Options{
			Addr:	  "localhost:6379",
			Password: "",
			DB:		  0,
		})
	} else {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(err)
		}
		c = redis.NewClient(opt)
	}

	if err := c.Ping(ctx).Err(); err != nil {
		panic("Unable to connect to redis " + err.Error())
	}
	client.c = c
	return client
}


//GetKey get key
func (client *redisClient) getKey(key string, src interface{}) error {
	val, err := client.c.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), &src)
	if err != nil {
		return err
	}
	return nil
}

//SetKey set key
func (client *redisClient) setKey(key string, value interface{}, expiration time.Duration) error {
	cacheEntry, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = client.c.Set(ctx, key, cacheEntry, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}


