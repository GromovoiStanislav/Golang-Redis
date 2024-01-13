package handlers

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"redis-example/internal/redisconn"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        var ctx = context.Background()

		red := redisconn.GetRedisConnection()

        res, err := red.Get(ctx, "user").Result()
        if err != nil {
            w.Write([]byte("Error!"))

        }
        w.Write([]byte(res))
    case "POST":
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            w.Write([]byte("Error!"))
        }
        defer r.Body.Close()
       
		var ctx = context.Background()

	
		red := redisconn.GetRedisConnection()
        _, err = red.SetNX(ctx,"user", body, 60*time.Second).Result()
        if err != nil {
            w.Write([]byte("Error!"))

        }
        w.Write([]byte(body))

    }

}