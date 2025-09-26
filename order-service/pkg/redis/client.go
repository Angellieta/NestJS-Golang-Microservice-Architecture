// order-service/pkg/redis/client.go
package redis

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// NewClient membuat dan mengembalikan koneksi klien Redis.
func NewClient() (*redis.Client, error) {
	// Mengambil URL Redis dari environment variable
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable not set")
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	// Ping untuk memastikan koneksi berhasil
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Redis connection successfully established")
	return client, nil
}