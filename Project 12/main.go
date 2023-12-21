package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

func main() {
	
	connectRedis(ctx)
	defer redisClient.Close()


    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("It is working 👊")
    })

	app.Get("/:id", verifyCache, func(c *fiber.Ctx) error {
        id := c.Params("id")
        res, err := http.Get("https://jsonplaceholder.typicode.com/users/" + id)
        if err != nil {
            return err
        }

        defer res.Body.Close()
        body, err := ioutil.ReadAll(res.Body)
        if err != nil {
            return err
        }

        cacheErr := redisClient.Set(ctx, id, body, 10*time.Second).Err()
        if cacheErr != nil {
            return cacheErr
        }

        data := toJson(body)
        return c.JSON(fiber.Map{"Data": data})
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

func verifyCache(c *fiber.Ctx) error {
    id := c.Params("id")
    val, err := redisClient.Get(ctx, id).Bytes()
    if err != nil {
        return c.Next()
    }

    data := toJson(val)
    return c.JSON(fiber.Map{"Cached": data})
}