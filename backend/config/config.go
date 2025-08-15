package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DatabaseURL string

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Server
	APIPort string
	Env     string

	// Authentication
	TokenExpiry int // hours

	// Supabase (if needed for additional features)
	SupabaseURL string
	SupabaseKey string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	// Parse token expiry with default
	tokenExpiry := 24 // default 24 hours
	if exp := os.Getenv("TOKEN_EXPIRY"); exp != "" {
		if parsed, err := strconv.Atoi(exp); err == nil {
			tokenExpiry = parsed
		}
	}

	// Parse Redis DB with default
	redisDB := 0 // default DB 0
	if db := os.Getenv("REDIS_DB"); db != "" {
		if parsed, err := strconv.Atoi(db); err == nil {
			redisDB = parsed
		}
	}

	cfg := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
		APIPort:       os.Getenv("API_PORT"),
		Env:           os.Getenv("ENV"),
		TokenExpiry:   tokenExpiry,
		SupabaseURL:   os.Getenv("SUPABASE_URL"),
		SupabaseKey:   os.Getenv("SUPABASE_KEY"),
	}
	return cfg, nil
}
