package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
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

	// Create the swagger handler.
	// We need to strip the /swagger prefix so the handler's internal routing works.
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)

	// Custom NotFound handler that is now more robust.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Safely check if the request path starts with the /swagger/ prefix.
		if strings.HasPrefix(r.URL.Path, "/swagger/") {
			// Strip the prefix and serve the swagger assets.
			http.StripPrefix("/swagger/", swaggerHandler).ServeHTTP(w, r)
			return
		}
		// Otherwise, serve the standard 404 response for our API.
		app.notFoundResponse(w, r)
	})

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
