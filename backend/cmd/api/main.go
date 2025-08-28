// File: backend/cmd/api/main.go

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cliffdoyle/internal/database"
	redisclient "github.com/cliffdoyle/internal/redis"
	"github.com/cliffdoyle/internal/repository"
	"github.com/cliffdoyle/internal/service"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
)

// @title           Hotel Management System API
// @version         1.0
// @description     This is the API for a MEWS-like hotel management system.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Type 'Bearer' followed by a space and a JWT token."

// config struct holds all configuration for the application.
type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
	redis struct {
		addr     string
		password string
		db       int
	}
}

// The `Services` and `Models` types are now collections of our other repository/service interfaces.
type Models struct {
	Permissions repository.PermissionRepository
	Users       repository.UserRepository
	Rooms       repository.RoomRepository
}
type Services struct {
	Users service.UserService
	Rooms  service.RoomService
}

// application struct holds the application-wide dependencies.
type application struct {
	config      config
	logger      *slog.Logger
	models      Models
	services    Services
	redis       *redis.Client
	db          *pgxpool.Pool
	metrics appMetrics // <-- ADD A CUSTOM REGISTRY
	metrics_reg  *prometheus.Registry
	// We will add models, services, repositories here later.
}

func main() {
	// Initialize a new structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("error loading .env file", "error", err)
		os.Exit(1)
	}

	// Initialize config struct
	var cfg config
	cfg.port = os.Getenv("PORT")
	if cfg.port == "" {
		cfg.port = "8080"
	}
	cfg.env = os.Getenv("ENV")
	if cfg.env == "" {
		cfg.env = "development"
	}
	cfg.db.dsn = os.Getenv("SUPABASE_URL")
	cfg.redis.addr = os.Getenv("REDIS_ADDR")

	//Establish database connection
	db, err := database.Connect(cfg.db.dsn)
	if err != nil {
		logger.Error("error connecting to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection established")

	// Establish Redis connection using the aliased package name
	redisClient, err := redisclient.Connect(cfg.redis.addr, cfg.redis.password, cfg.redis.db)
	if err != nil {
		logger.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}
	logger.Info("redis connection pool established")

	// --- Initialize Prometheus Registry ---
	metrics_reg := prometheus.NewRegistry()

	appMetrics := newMetrics(metrics_reg)

	// --- Initialize repositories ---
	userRepo := repository.NewUserRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	roomRepo := repository.NewRoomRepository(db)

	// --- Initialize services ---
	userService := service.NewUserService(userRepo, redisClient)
	roomService := service.NewRoomService(roomRepo) 

	// Initialize application struct
	app := &application{
		config: cfg,
		logger: logger,
		services: Services{
			Users: userService,
			Rooms: roomService,
		},
		models: Models{
			Users:       userRepo,
			Permissions: permissionRepo,
			Rooms:       roomRepo,
		},
		redis: redisClient,
		db:    db,
		metrics: appMetrics,
		metrics_reg: metrics_reg,
	}

	// Create a new server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.port),
		Handler:      app.routes(), // Our new routes function
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Graceful shutdown logic
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	// Start the server
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Error("server startup failed", "error", err)
		os.Exit(1)
	}

	// Wait for shutdown to complete
	if err = <-shutdownError; err != nil {
		logger.Error("server shutdown failed", "error", err)
	} else {
		logger.Info("server stopped gracefully")
	}
}
