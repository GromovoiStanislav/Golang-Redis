## Redis and Gin leaderboard API

In this project, we will explore go-redis and use its pipeline feature to build a leaderboard API. The API will use Gin and Redis' sorted sets under the hood. It will expose the following endpoints:

- GET /points/:username — to get a user's score and their rank in the overall leaderboard
- POST /points — to add or update a user and their score. This endpoint will also return the new rank of the user
- GET /leaderboard — returns the current leaderboard, with users sorted in ascending order of their ranks

#### Setup:

```
go mod init redis-example
go get github.com/redis/go-redis/v9
go get github.com/joho/godotenv
go get -u github.com/gin-gonic/gin
```

#### Start:

```
go run main.go
```
