// File: backend/internal/redis/client.go

package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Connect establishes a connection to the Redis server.
func Connect(addr, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
