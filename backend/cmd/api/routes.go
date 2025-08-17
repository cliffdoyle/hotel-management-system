package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, i interface{}) {
		app.serverErrorResponse(w, r, fmt.Errorf("panic recovered: %v", i))
	}

	// --- Public User Routes ---
	router.HandlerFunc(http.MethodPost, "/v1/users/register", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/login", app.loginHandler)

	// --- Other Public Routes ---
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/refresh", app.refreshTokenHandler)

	// Add protected profile management routes in the future here
	
	return app.authenticate(router)
}