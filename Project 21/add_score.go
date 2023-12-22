package main

 import (
     "context"
	 
     "github.com/redis/go-redis/v9"
 )  

 func addScore(c *redis.Client, p map[string]interface{}) (map[string]interface{}, error) {

     ctx := context.TODO()

     nickname  := p["nickname"].(string)
     steps     := p["steps"].(float64)

     //Validate data here in a production environment

     err := c.ZAdd(ctx, "app_users", redis.Z{
             Score:  steps,
             Member: nickname,
         }).Err()

     if err != nil {
         return nil, err
     } 

     rank := c.ZRank(ctx, "app_users", p["nickname"].(string))

     if err != nil {
         return nil, err
     } 

     response := map[string]interface{}{
                     "data": map[string]interface{}{
                         "nickname": p["nickname"].(string),
                         "rank":     rank.Val(),
                      },
                 }

     return response, nil
 }