package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	http.HandleFunc("/scores", httpHandler)        
	http.ListenAndServe(":3000", nil)
}



func httpHandler(w http.ResponseWriter, req *http.Request) {            
	var err error

	err = godotenv.Load()
	if err != nil {
		panic(err)
	}
	redisURL := os.Getenv("REDIS_URL")


	var redisClient *redis.Client
	if redisURL == "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:	  "localhost:6379",
			Password: "",
			DB:		  0,
		})
	} else {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(err)
		}
		redisClient = redis.NewClient(opt)
	}
	defer redisClient.Close()
	
	
	params := map[string]interface{}{}

	resp := map[string]interface{}{}

	if req.Method == "GET" {
		for k, v := range req.URL.Query() {
			params[k] = v[0]
		}    
		resp, err = getScores(redisClient, params)

	} else if req.Method == "POST" {
		err = json.NewDecoder(req.Body).Decode(&params)
		resp, err = addScore(redisClient, params)
		
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	if err != nil {
		resp = map[string]interface{}{
				   "error": err.Error(),
			   }
	} else {

		if encodingErr := enc.Encode(resp); encodingErr != nil {
			fmt.Println("{ error: " + encodingErr.Error() + "}")
		}   
	}

}