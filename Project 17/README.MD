## Real Time Chat with Gorilla/WebSocket and Redis

#### Setup:

```
go mod init redis-example
go get github.com/redis/go-redis/v9
go get github.com/gorilla/websocket
go get github.com/joho/godotenv
```

#### Start server:

```
go run server.go
```

#### Start clients:

```
go run client.go
```
