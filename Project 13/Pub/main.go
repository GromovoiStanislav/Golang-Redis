package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

var redisClient *redis.Client
var ctx = context.Background()

func main() {
	
	connectRedis(ctx)
	defer redisClient.Close()


    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("It is working 👊")
    })

	app.Post("/", func(c *fiber.Ctx) error {
        user := new(User)

        if err := c.BodyParser(user); err != nil {
            panic(err)
        }

        payload, err := json.Marshal(user)
        if err != nil {
            panic(err)
        }

        if err := redisClient.Publish(ctx, "send-user-data", payload).Err(); err != nil {
            panic(err)
        }

        return c.SendStatus(200)
    })

    app.Listen(":3000")
}

func connectRedis(ctx context.Context) {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	redisURL := os.Getenv("REDIS_URL")


	if redisURL == "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:	  "localhost:6379",
			Password: "", // no password set
			DB:		  0,  // use default DB
		})
	} else {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(err)
		}
		redisClient = redis.NewClient(opt)
	}
}

