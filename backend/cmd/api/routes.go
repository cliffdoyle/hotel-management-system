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

	// Our first endpoint: a health check.
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	return router
}
