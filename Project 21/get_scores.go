package main

 import (
     "context"
     "fmt"
     "strconv"
	 
     "github.com/redis/go-redis/v9"

 )  

 func getScores(c *redis.Client, p map[string]interface{}) (map[string]interface{}, error) {

     ctx := context.TODO()

     start, err := strconv.ParseInt(fmt.Sprint(p["start"]), 10, 64)

     if err != nil {
         return nil, err
     }

     stop, err  := strconv.ParseInt(fmt.Sprint(p["stop"]), 10, 64)

     if err != nil {
         return nil, err
     }

     total, err  := c.ZCount(ctx, "app_users", "-inf", "+inf").Result() //int64     

     if err != nil {
         return nil, err
     }

     scores, err := c.ZRevRangeWithScores(ctx, "app_users", start, stop).Result() //highest to lowest score

     if err != nil {
         return nil, err
     }

     data   := []map[string]interface{}{}            

     for _, z := range scores {

         record := map[string]interface{}{}
         rank   := c.ZRank(ctx, "app_users", z.Member.(string))

         if err != nil {
             return nil, err
         } 

         record["nickname"] = z.Member.(string)
         record["score"]    = z.Score 
         record["rank"]     = rank.Val()

         data = append(data, record)   

     }

     countPerRequest := stop - start + 1

     if stop == -1 {
         countPerRequest = total
     } 

     response := map[string]interface{}{
                     "data": data,
                     "meta": map[string]interface{}{
                         "start": start,
                         "stop":  stop,
                         "per_request": countPerRequest,
                         "total": total,
                     },
                 }

     return response, nil

 }