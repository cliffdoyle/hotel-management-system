package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cliffdoyle/config"
	"github.com/cliffdoyle/db"
	"github.com/go-chi/cors"
	"github.com/julienschmidt/httprouter"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{"status": "ok"}

	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "failed to send response to client", http.StatusInternalServerError)
	}

	w.Write(res)

	fmt.Println(res)
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if _, err := db.ConnectPostgres(cfg); err != nil {
		log.Fatalf("Failed to connect to Supabase: %v", err)
	}
	if _, err := db.ConnectRedis(cfg); err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}

	router := httprouter.New()

	router.GET("/health", healthCheckHandler)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Allow our Vite frontend
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})
	handler := c.Handler(router)

	serverAddr := ":" + cfg.APIPort
	log.Printf("Starting server on http://localhost%s", serverAddr)

	if err := http.ListenAndServe(serverAddr, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
