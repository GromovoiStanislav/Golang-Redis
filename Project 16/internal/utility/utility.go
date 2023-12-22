package utility

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Failed to read .env")
	}
}

func getEnv(key string) string {
	loadEnv()
	value, status := os.LookupEnv(key)
	if !status {
		log.Fatalf("Missing environment variable %s\n", key)
	}
	return value
}

func Address() string {
	return getEnv("REDIS_ADDRESS")
}

func Username() string {
	return getEnv("REDIS_USER")
}

func Password() string {
	return getEnv("REDIS_PASSWORD")
}

func Database() int {
	databaseStr := getEnv("REDIS_DATABASE")
	database, err := strconv.Atoi(databaseStr)
	if err != nil {
		log.Fatalln("Database environment variable must be a number")
	}
	return database
}