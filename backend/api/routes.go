package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// setupRoutes configures all application routes
func (app *Application) setupRoutes() {
	// Health check endpoint
	app.Router.GET("/health", app.healthCheckHandler)

	// API v1 routes
	v1 := "/api/v1"

	// Authentication routes
	app.Router.POST(v1+"/auth/register", app.registerHandler)
	app.Router.POST(v1+"/auth/login", app.loginHandler)
	app.Router.POST(v1+"/auth/logout", app.requireAuth(app.logoutHandler))

	// User management routes
	app.Router.GET(v1+"/users/profile", app.requireAuth(app.getUserProfileHandler))
	app.Router.PUT(v1+"/users/profile", app.requireAuth(app.updateUserProfileHandler))
	app.Router.PUT(v1+"/users/password", app.requireAuth(app.changePasswordHandler))

	// Protected routes will be added here as we build more features
}

// setupMiddleware configures global middleware
func (app *Application) setupMiddleware() {
	// Add CORS middleware
	app.Router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})
}

// requireAuth is a middleware wrapper for protected routes
func (app *Application) requireAuth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Authentication logic will be implemented here
		// For now, just call the next handler
		next(w, r, ps)
	}
}
