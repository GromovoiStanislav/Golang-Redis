## Caching using Redis and Fiber

#### Setup:

```
go mod init redis-examples
go get github.com/redis/go-redis/v9
go get github.com/joho/godotenv
go get github.com/gofiber/fiber/v2
```

#### Start:

```
go run subscriber.go
go run publisher.go
```
