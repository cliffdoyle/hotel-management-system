// File: backend/cmd/api/routes.go

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Convert our custom error handlers to http.HandlerFunc
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// A handler for catching panics and recovering.
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, i interface{}) {
		app.logger.Error("recovered from panic", "panic", i)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// Create a dummy user handler for testing the middleware.
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.requireAuthenticatedUser(app.getUserHandler)) // Dummy handler

	// A sample route with permission check.
	router.HandlerFunc(http.MethodPost, "/v1/users", app.requirePermission("users:write", app.createUserHandler)) // Dummy handler

	// Our first endpoint: a health check.
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	return app.authenticate(router)
}

// Dummy handler functions to test our middleware
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
}
func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusCreated, envelope{"message": "user created successfully (not really)"}, nil)
}
