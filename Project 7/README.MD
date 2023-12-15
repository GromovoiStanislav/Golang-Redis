## Redis Stream example

#### Setup:

```
go mod init redis-examples
go get github.com/redis/go-redis/v9
go get github.com/joho/godotenv
```

#### Start:

```
go run producer.go
go run сonsumer_1.go
go run сonsumer_2.go
go run сonsumer_3.go
go run reader.go
```
