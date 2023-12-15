package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"redis-example/internal"
	internal_redis "redis-example/internal/redis"
)


func main() {
	errC, err := run()
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run() (<-chan error, error) {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	// Чтение переменной окружения REDIS_URL
	redisURL := os.Getenv("REDIS_URL")


	rdb:= internal_redis.NewRedis(redisURL)
	//defer rdb.Close()


	srv := &Server{
		rdb:    rdb,
		done:   make(chan struct{}),
	}


	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)


	go func() {
		<-ctx.Done()

		log.Println("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer func() {
			rdb.Close()
			stop()
			cancel()
			close(errC)
		}()

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		log.Println("Shutdown completed")

	}()

	
	go func() {
		log.Println("Listening and serving")

		if err := srv.ListenAndServe(); err != nil {
			errC <- err
		}
	}()

	return errC, nil
}

type Server struct {
	rdb    *redis.Client
	pubsub *redis.PubSub
	done   chan struct{}
}

func (s *Server) ListenAndServe() error {
	pubsub := s.rdb.PSubscribe(context.Background(), "tasks.*") // Pattern-matching subscription

	_, err := pubsub.Receive(context.Background())
	if err != nil {
		return fmt.Errorf("pubsub.Receive %w", err)
	}

	s.pubsub = pubsub

	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			log.Printf(fmt.Sprintf("Received message: %s", msg.Channel))

			switch msg.Channel {
			case "tasks.event.updated", "tasks.event.created":
				var task internal.Task

				if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&task); err != nil {
					log.Println("Ignoring message, invalid")
					continue
				}

				log.Println(task)
				// ...

			case "tasks.event.deleted":
				var id string

				if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&id); err != nil {
					log.Println("Ignoring message, invalid")
					continue
				}

				log.Printf("ID to delete:  %s",id)
				// ...
			}
		}

		// ...

		log.Println("No more messages to consume. Exiting.")

		s.done <- struct{}{}
	}()

	return nil
}

// Shutdown ...
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server")

	s.pubsub.Close()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context.Done: %w", ctx.Err())

		case <-s.done:
			return nil
		}
	}
}