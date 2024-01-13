package main

import (
    "log"
    "net/http"

    "redis-example/internal/handlers"
)

func main() {
    port := ":3000"
    http.HandleFunc("/users", handlers.UsersHandler)
    log.Fatal(http.ListenAndServe(port, nil))
}