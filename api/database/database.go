package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

// Create a new Redis client to connect to the Redis server
func CreateClient(db_no int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "db" + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db_no,
	})

	return rdb
}
