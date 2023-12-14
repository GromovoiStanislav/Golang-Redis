package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)


var client  *redis.Client

var ctx = context.Background()



func TestMain(m *testing.M) {
	// Вызываем нашу функцию инициализации перед выполнением всех тестов
	setup()

	// Запускаем тесты
	exitCode := m.Run()

	// Завершаем ресурсы после выполнения всех тестов
	teardown()

	// Возвращаем exit code
	os.Exit(exitCode)
}

func setup() {
	// Инициализация ресурсов перед запуском тестов, например, соединение с базой данных
		
	// Загрузка переменных окружения из файла .env
	godotenv.Load()
	// Чтение переменной окружения REDIS_URL
	redisURL := os.Getenv("REDIS_URL")

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



}

func teardown() {
	// Завершение ресурсов после выполнения всех тестов, например, закрытие соединения с базой данных
	if client != nil {
		client.Close()
	}
}



func TestConnection(t *testing.T) {
	assert.NotNil(t, client)
}

func TestPing(t *testing.T) {
	result, err := client.Ping(ctx).Result()
	assert.Nil(t, err)
	assert.Equal(t, "PONG", result)
}

func TestString(t *testing.T) {
	client.SetEx(ctx, "name", "Eko Kurniawan", 3*time.Second)

	result, err := client.Get(ctx, "name").Result()
	assert.Nil(t, err)
	assert.Equal(t, "Eko Kurniawan", result)

	time.Sleep(5 * time.Second)

	result, err = client.Get(ctx, "name").Result()
	assert.NotNil(t, err)
	assert.Equal(t, redis.Nil, err)
}

func TestList(t *testing.T) {
	client.RPush(ctx, "names", "Eko")
	client.RPush(ctx, "names", "Kurniawan")
	client.RPush(ctx, "names", "Khannedy")

	assert.Equal(t, "Eko", client.LPop(ctx, "names").Val())
	assert.Equal(t, "Kurniawan", client.LPop(ctx, "names").Val())
	assert.Equal(t, "Khannedy", client.LPop(ctx, "names").Val())
}

func TestSet(t *testing.T) {
	client.SAdd(ctx, "students", "Eko")
	client.SAdd(ctx, "students", "Eko")
	client.SAdd(ctx, "students", "Kurniawan")
	client.SAdd(ctx, "students", "Kurniawan")
	client.SAdd(ctx, "students", "Khannedy")
	client.SAdd(ctx, "students", "Khannedy")

	assert.Equal(t, int64(3), client.SCard(ctx, "students").Val())
	assert.Equal(t, []string{"Eko", "Kurniawan", "Khannedy"}, client.SMembers(ctx, "students").Val())

	client.Del(ctx, "students")
}

func TestSortedSet(t *testing.T) {
	client.ZAdd(ctx, "scores", redis.Z{Score: 100, Member: "Eko"})
	client.ZAdd(ctx, "scores", redis.Z{Score: 85, Member: "Budi"})
	client.ZAdd(ctx, "scores", redis.Z{Score: 95, Member: "Joko"})
	client.ZAdd(ctx, "scores", redis.Z{Score: 75, Member: "Loo"})

	assert.Equal(t, []string{"Loo", "Budi", "Joko", "Eko"}, client.ZRange(ctx, "scores", 0, -1).Val())

	assert.Equal(t, "Loo", client.ZPopMin(ctx, "scores").Val()[0].Member)
	assert.Equal(t, "Eko", client.ZPopMax(ctx, "scores").Val()[0].Member)
	assert.Equal(t, "Joko", client.ZPopMax(ctx, "scores").Val()[0].Member)
	assert.Equal(t, "Budi", client.ZPopMax(ctx, "scores").Val()[0].Member)
}

func TestHash(t *testing.T) {
	client.HSet(ctx, "user:1", "id", "1")
	client.HSet(ctx, "user:1", "name", "Eko")
	client.HSet(ctx, "user:1", "email", "eko@example.com")

	user := client.HGetAll(ctx, "user:1").Val()

	assert.Equal(t, "1", user["id"])
	assert.Equal(t, "Eko", user["name"])
	assert.Equal(t, "eko@example.com", user["email"])

	client.Del(ctx, "user:1")
}

func TestGeoPoint(t *testing.T) {
	client.GeoAdd(ctx, "sellers", &redis.GeoLocation{
		Name:      "Toko A",
		Longitude: 106.818489,
		Latitude:  -6.178966,
	})
	client.GeoAdd(ctx, "sellers", &redis.GeoLocation{
		Name:      "Toko B",
		Longitude: 106.821568,
		Latitude:  -6.180662,
	})

	distance := client.GeoDist(ctx, "sellers", "Toko A", "Toko B", "km").Val()
	assert.Equal(t, 0.3892, distance)

	sellers := client.GeoSearch(ctx, "sellers", &redis.GeoSearchQuery{
		Longitude:  106.819143,
		Latitude:   -6.180182,
		Radius:     5,
		RadiusUnit: "km",
	}).Val()

	assert.Equal(t, []string{"Toko A", "Toko B"}, sellers)

	client.Del(ctx, "sellers")
}

func TestHyperLogLog(t *testing.T) {
	client.PFAdd(ctx, "visitors", "eko", "kurniawan", "khannedy")
	client.PFAdd(ctx, "visitors", "eko", "budi", "joko")
	client.PFAdd(ctx, "visitors", "rully", "budi", "joko")

	total := client.PFCount(ctx, "visitors").Val()
	assert.Equal(t, int64(6), total)

	client.Del(ctx, "visitors")
}

func TestPipeline(t *testing.T) {
	_, err := client.Pipelined(ctx, func(pipeliner redis.Pipeliner) error {
		pipeliner.SetEx(ctx, "name", "Eko", 5*time.Second)
		pipeliner.SetEx(ctx, "address", "Indonesia", 5*time.Second)
		return nil
	})
	assert.Nil(t, err)

	assert.Equal(t, "Eko", client.Get(ctx, "name").Val())
	assert.Equal(t, "Indonesia", client.Get(ctx, "address").Val())
}

func TestTransaction(t *testing.T) {
	_, err := client.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		pipeliner.SetEx(ctx, "name", "Joko", 5*time.Second)
		pipeliner.SetEx(ctx, "address", "Cirebon", 5*time.Second)
		return nil
	})
	assert.Nil(t, err)

	assert.Equal(t, "Joko", client.Get(ctx, "name").Val())
	assert.Equal(t, "Cirebon", client.Get(ctx, "address").Val())
}

func TestCreateConsumerGroup(t *testing.T) {
	client.XGroupCreate(ctx, "members", "group-1", "0")
	client.XGroupCreateConsumer(ctx, "members", "group-1", "consumer-1")
	client.XGroupCreateConsumer(ctx, "members", "group-1", "consumer-2")
}

func TestSendMessage(t *testing.T) {
    message := map[string]interface{}{"key1": "value1", "key2": "value2"}

    // Отправляем сообщение в поток "members"
    client.XAdd(ctx, &redis.XAddArgs{
        Stream: "members",
        Values: message,
    })

}

func TestConsumeStream(t *testing.T) {
	streams := client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    "group-1",
		Consumer: "consumer-1",
		Streams:  []string{"members", ">"},
		Count:    2,
		Block:    10 * time.Second,
	}).Val()

	for _, stream := range streams {
		for _, message := range stream.Messages {
			fmt.Println(message.ID)
			fmt.Println(message.Values)
		}
	}
}
