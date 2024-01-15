package main

import (
	"context"
	"fmt"
	"redis-example/internal/redisconn"
	"time"
)

func main() {
   
    {
        var ctx = context.Background()
        red := redisconn.GetRedisConnection()
        red.SetNX(ctx, "user", "user-json", 60*time.Second).Result()
    }


    {
        var ctx = context.Background()
        red := redisconn.GetRedisConnection()
        res, _ := red.Get(ctx, "user").Result()
        fmt.Println(res)
    }

}