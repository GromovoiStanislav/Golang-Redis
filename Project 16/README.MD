## Using Redis in Go simple examples

This repo contains various examples of using Redis in Go. You'll learn how to perform the following actions:

- Connect to a Redis server from your Go program.
- Add data to a Redis database.
- Read previously cached data in the database.
- Update a cached value.
- Delete a cached value.

#### Setup:

```
go mod init redis-example
go get github.com/redis/go-redis/v9
go get github.com/joho/godotenv
```

#### Start:

```bash
# Connecting to the Redis server
go run cmd/connect/connect.go

# Adding data to the cache
go run cmd/set/set.go

# Reading data from the cache
go run cmd/get/get.go

# Updating data in the cache
go run cmd/update/update.go

# Deleting data from the cache
go run cmd/delete/delete.go

# Putting it all together
go run cmd/demo/demo.go
```
