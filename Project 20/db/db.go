package db

import (
	"context"
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	Client *redis.Client
}

var (
	ErrNil = errors.New("no matching record found in redis database")
	Ctx    = context.TODO()
)

func NewDatabase() (*Database, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	redisURL := os.Getenv("REDIS_URL")


	var redisClient *redis.Client

	if redisURL == "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:	  "localhost:6379",
			Password: "", // no password set
			DB:		  0,  // use default DB
		})
		
	} else {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			return nil, err
		}
		redisClient = redis.NewClient(opt)
	}


	if err := redisClient.Ping(Ctx).Err(); err != nil {
		return nil, err
	}

	return &Database{
		Client: redisClient,
	}, nil
}

