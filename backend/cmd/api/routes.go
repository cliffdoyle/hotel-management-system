package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// The `routes` function now returns a fully wrapped http.Handler.
func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, i interface{}) {
		app.serverErrorResponse(w, r, fmt.Errorf("panic recovered: %v", i))
	}
	
	// --- Public Routes ---
	router.HandlerFunc(http.MethodPost, "/v1/users/register", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/login", app.loginHandler)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/refresh", app.refreshTokenHandler)
	router.Handler(http.MethodGet, "/metrics", app.metricsHandler())

	// Chain our middleware together in the correct order.
	// Outermost middleware is applied first.
	var chain http.Handler
	chain = router
	chain = app.withMetrics(chain)
	chain = app.withRateLimit(chain)
	chain = app.authenticate(chain) // Authenticate is after rate limiting
	chain = app.withCORS(chain)
	chain = app.withSecurityHeaders(chain)
	chain = app.withLogging(chain)
	chain = app.withRequestID(chain)
	
	return chain
}