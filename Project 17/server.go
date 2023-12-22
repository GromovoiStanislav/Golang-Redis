package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var client *redis.Client
var Users map[string]*websocket.Conn
var sub *redis.PubSub
var upgrader = websocket.Upgrader{}

const chatChannel = "chats"

func init() {
	Users = map[string]*websocket.Conn{}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	// Подключение к Redis (замените параметры подключения на свои)
	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Проверка подключения к Redis
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis:", pong)

	broadcast()

	http.HandleFunc("/chat/", chat)
	server := http.Server{Addr: ":8080", Handler: nil}

	// Запуск сервера
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Server error:", err)
		}
	}()



	// ... Закрытие соединений, отписка и остановка сервера
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	<-exit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Очистка ресурсов
	for _, conn := range Users {
		if err := conn.Close(); err != nil {
			log.Println("Error closing connection:", err)
		}
	}

	if err := sub.Close(); err != nil {
		log.Println("Error unsubscribing from Redis channel:", err)
	}

	if err := client.Close(); err != nil {
		log.Println("Error closing Redis client:", err)
	}

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}
}

func chat(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimPrefix(r.URL.Path, "/chat/")

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// 1. Создание веб-сокет соединения
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func() {
		// Закрытие соединения при выходе из функции
		if err := c.Close(); err != nil {
			log.Println("Error closing connection:", err)
		}
	}()

	// 2. Ассоциирование пользователя (имени) с фактическим соединением
	Users[user] = c
	fmt.Println(user, "in chat")

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break // Выход из цикла при ошибке чтения (например, если клиент отключился)
		}

		// 3. При поступлении сообщения через соединение, публикация сообщения в канал Redis
		client.Publish(context.Background(), chatChannel, user+":"+string(message))
	}
}


func broadcast() {
	go func() {
		sub = client.Subscribe(context.Background(), chatChannel)
		messages := sub.Channel()
		for message := range messages {
			from := strings.Split(message.Payload, ":")[0]
			// 3. Если сообщение получено на канале Redis, рассылка его всем подключенным сессиям (пользователям)
			for user, conn := range Users {
				if from != user {
					conn.WriteMessage(websocket.TextMessage, []byte(message.Payload))
				}
			}
		}
	}()
}
