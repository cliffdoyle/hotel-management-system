package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/cliffdoyle/config"
	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

// Application struct holds all dependencies for dependency injection
type Application struct {
	Config *config.Config
	DB     *sql.DB
	Redis  *redis.Client
	Router *httprouter.Router
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test Redis connection
	if err := redisClient.Ping(redisClient.Context()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Initialize router
	router := httprouter.New()

	// Create application instance with all dependencies
	app := &Application{
		Config: cfg,
		DB:     db,
		Redis:  redisClient,
		Router: router,
	}

	// Setup routes
	app.setupRoutes()

	// Setup middleware
	app.setupMiddleware()

	log.Printf("Server starting on port %s", cfg.APIPort)
	log.Fatal(http.ListenAndServe(":"+cfg.APIPort, app.Router))
}
