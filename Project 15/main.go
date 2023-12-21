package main

import (
	"context"
	"fmt"
	"log"
	"os"


	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)


var redisClient *redis.Client

var ctx = context.Background()


func main() {
	// Connecting to a Redis Instance
	redisClient = createRedisClient()
	defer redisClient.Close()

	
	response := redisClient.Ping(ctx).Val() 
	fmt.Println(response) // PONG


	// Inserting Values Into Your Redis Instance
	fmt.Println(insertIt("key1","value1")) // insertion successful


	// Querying and Reading from your Redis instance
	fmt.Println("key1", getIt("key1"))


	// Deleting Key-Value Pairs from your Redis Instance
	fmt.Println("deleted:", deleteIt("key1")) // deleted: 1


	// Querying and Reading from your Redis instance
	fmt.Println("key1", getIt("key1")) // key1 does not exist | key1 <nil>  


	// Working With Redis Lists
	insertResult, err := redisClient.LPush(ctx,"Companies", "Google", "Microsoft", "Netflix", "Amazon", "Uber", "Airbnb", "Microsoft").Result() 
	if err != nil { 
		fmt.Println(err) 
	} else {
		fmt.Println("inserted:", insertResult) // inserted: 7
	}

	rangeResult, err := redisClient.LRange(ctx, "Companies", 1, 3).Result() 
	if err != nil { 
		fmt.Println(err) 
	} else {
		fmt.Println("range result:", rangeResult) // [Airbnb Uber Amazon]
	}

	result, err := redisClient.LRem(ctx, "Companies", 2, "Microsoft").Result()  
	if err != nil { 
		fmt.Println(err) 
	} else {
		fmt.Println("removed:", result) // removed: 2
	}

	data, err := redisClient.LPop(ctx, "Companies").Result() 
	if err != nil { 
		fmt.Println(err) 
	} else {
		fmt.Println("Companie:", data) // Companie: Airbnb
	}



	// Printing All Keys in a Redis Cluster
	allkeys, err := redisClient.Keys(ctx,"*").Result() 
	if err != nil { 
		fmt.Println(err) 
	} else {
		fmt.Println("all keys:", allkeys) // all keys: [Companies] 
	}


	
	// Using Go-Redis Unsupported Commands
	unsupported()// key does not exists
}


// Connecting to a Redis Instance
func createRedisClient() *redis.Client { 

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	redisURL := os.Getenv("REDIS_URL")

	if redisURL == "" {
		redisClient := redis.NewClient(&redis.Options{
			Addr:	  "localhost:6379",
			Password: "", // no password set
			DB:		  0,  // use default DB
		})
		return redisClient

	} else {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(err)
		}
		redisClient := redis.NewClient(opt)
		return redisClient
	}
} 



// Inserting Values Into Your Redis Instance
func insertIt(key, value string) any { 
	err := redisClient.Set(ctx, key, value, 0).Err() 
	if err != nil { 
		return err 
	} 
	return "insertion successful" 
 } 

 // Querying and Reading from your Redis instancefunc insertIt(key, value string) any { 
func getIt(key string) any { 
	val, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
        log.Printf("%v does not exist", key)
		return nil
	} else if err != nil {
		log.Printf("Ошибка загрузки ключа %v:%v ", key, err)
		return nil
	}
	return val 
} 

// Deleting Key-Value Pairs from your Redis Instance
func deleteIt(key string) any { 
    deleted, err := redisClient.Del(ctx, key).Result() 
    if err != nil { 
         log.Println(err) 
    } 
    return deleted 
}

func unsupported() { 
	result, err := redisClient.Do(ctx, "get", "key1").Result() 
	if err != nil { 
		if err == redis.Nil { 
			log.Println("key does not exists") 
			return 
		} 
		panic(err) 
	} 
	fmt.Println(result.(string)) 
}