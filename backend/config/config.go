package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr   string
	SupabaseURL string
	APIPort     string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg:=&Config{
		RedisAddr: os.Getenv("REDIS_ADDR"),
		SupabaseURL: os.Getenv("SUPABASE_URL"),
		APIPort: os.Getenv("API_PORT"),
	}
	return cfg,nil
}
