## Redis examples

#### Setup:

```
go mod init redis-examples
go get github.com/redis/go-redis/v9
go get github.com/joho/godotenv
```

#### Start:

```bash
go run example_1.go     # String
go run example_2.go     # Hash
go run example_3.go     # list
go run example_4.go     # Pipeline
go run example_5.go     # HyperLogLog
go run example_6.go     # Incr
go run example_7.go     # Delete keys without a ttl
go run example_8.go     # Delete keys without a ttl
go run example_9.go     # PubSub
go run example_10.go    # Bloom & Cuckoo filters
```
