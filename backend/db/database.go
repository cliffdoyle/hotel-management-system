package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/cliffdoyle/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("successfully connected to Redis")
	return rdb, nil
}

func ConnectPostgres(cfg *config.Config) (*sql.DB, error) {
	pool, err := sql.Open("Postgres", cfg.SupabaseURL)
	if err != nil {
		return nil, err
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()

	err = pool.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Successfully connected to PostgreSQL")
	return pool, nil
}
