package cache

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client    *redis.Client
	initOnce  sync.Once
	initError error
)

// InitRedis initializes a Redis client using environment configuration and verifies connectivity.
func InitRedis() error {
	initOnce.Do(func() {
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			addr = "localhost:6379"
		}

		password := os.Getenv("REDIS_PASSWORD")

		client = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			initError = fmt.Errorf("redis ping failed: %w", err)
		}
	})

	return initError
}

// Client returns the initialized Redis client instance.
func Client() *redis.Client {
	return client
}

// Close releases the Redis client resources when shutting down the application.
func Close() error {
	if client == nil {
		return nil
	}
	return client.Close()
}
