package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// User is a struct representing newly registered users
type User struct {
	Username string
	Email    string
}

// MarshalBinary encodes the struct into a binary blob
// Here I cheat and use regular json :)
func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// UnmarshalBinary decodes the struct into a User
func (u *User) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	return nil
}

// Names Some Non-Random name lists used to generate Random Users
var Names []string = []string{"Jasper", "Johan", "Edward", "Niel", "Percy", "Adam", "Grape", "Sam", "Redis", "Jennifer", "Jessica", "Angelica", "Amber", "Watch"}

// SirNames Some Non-Random name lists used to generate Random Users
var SirNames []string = []string{"Ericsson", "Redisson", "Edisson", "Tesla", "Bolmer", "Andersson", "Sword", "Fish", "Coder"}

// EmailProviders Some Non-Random email lists used to generate Random Users
var EmailProviders []string = []string{"Hotmail.com", "Gmail.com", "Awesomeness.com", "Redis.com"}

func main() {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	// Чтение переменной окружения REDIS_URL
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
			panic(err)
		}
		redisClient = redis.NewClient(opt)
	}
	defer redisClient.Close()

	// Generate a new background context that  we will use
	ctx := context.Background()
	

	// Ping the Redis server and check if any errors occured
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		// Sleep for 3 seconds and wait for Redis to initialize
		time.Sleep(3 * time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}
	}
		

	// Loop and randomly generate users on a random timer
	for {
		// Publish a generated user to the new_users channel
		err := redisClient.Publish(ctx, "new_users", GenerateRandomUser()).Err()
		if err != nil {
			panic(err)
		}
		// Sleep random time
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(4)
		time.Sleep(time.Duration(n) * time.Second)
	}
}


// GenerateRandomUser creates a random user, dont care too much about this.
func GenerateRandomUser() *User {
	rand.Seed(time.Now().UnixNano())
	nameMax := len(Names)
	sirNameMax := len(SirNames)
	emailProviderMax := len(EmailProviders)

	nameIndex := rand.Intn(nameMax-1) + 1
	sirNameIndex := rand.Intn(sirNameMax-1) + 1
	emailIndex := rand.Intn(emailProviderMax-1) + 1

	return &User{
		Username: Names[nameIndex] + " " + SirNames[sirNameIndex],
		Email:    Names[nameIndex] + SirNames[sirNameIndex] + "@" + EmailProviders[emailIndex],
	}
}