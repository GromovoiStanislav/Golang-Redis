package redis

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"

	"redis-example/internal"
)



func NewRedis(redisURL string) *redis.Client {
	var client *redis.Client
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
	return client
}


// Task represents the repository used for publishing Task records.
type Task struct {
	client *redis.Client
}

// NewTask instantiates the Task repository.
func NewTask(client *redis.Client) *Task {
	return &Task{
		client: client,
	}
}

// Created publishes a message indicating a task was created.
func (t *Task) Created(ctx context.Context, task internal.Task) error {
	return t.publish(ctx, "Task.Created", "tasks.event.created", task)
}

// Deleted publishes a message indicating a task was deleted.
func (t *Task) Deleted(ctx context.Context, id string) error {
	return t.publish(ctx, "Task.Deleted", "tasks.event.deleted", id)
}

// Updated publishes a message indicating a task was updated.
func (t *Task) Updated(ctx context.Context, task internal.Task) error {
	return t.publish(ctx, "Task.Updated", "tasks.event.updated", task)
}

func (t *Task) publish(ctx context.Context, spanName, channel string, e interface{}) error {
	var b bytes.Buffer

	if err := json.NewEncoder(&b).Encode(e); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.Encode")
	}

	res := t.client.Publish(ctx, channel, b.Bytes())
	if err := res.Err(); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "client.Publish")
	}

	return nil
}