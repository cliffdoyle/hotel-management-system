// File: backend/cmd/api/middleware.go

package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/cliffdoyle/internal/models"
	"github.com/cliffdoyle/internal/repository"
	"github.com/cliffdoyle/internal/tokens"
	"github.com/google/uuid"
)

// A custom contextKey type to avoid key collisions in context.
type contextKey string

const userContextKey = contextKey("user")

// responseWriter is a custom response writer that captures the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// // withMetrics records Prometheus metrics for each request.
// func (app *application) withMetrics(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Use our custom response writer again.
// 		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

// 		// Start the timer.
// 		start := time.Now()

// 		// Let the request proceed.
// 		next.ServeHTTP(rw, r)

// 		duration := time.Since(start).Seconds()

// 		// Record the metrics.
// 		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rw.statusCode)).Inc()
// 		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
// 	})
// }

// withRequestID generates a unique ID for each request.
func (app *application) withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "requestID", requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// withLogging logs details about each request and its response.
func (app *application) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Use our custom response writer to capture status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Let the request proceed down the chain
		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		requestID, _ := r.Context().Value("requestID").(string)

		app.logger.Info("request processed",
			"id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"status", rw.statusCode,
			"duration", duration.String(),
		)
	})
}

// withSecurityHeaders adds common security headers to every response.
func (app *application) withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline';")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		next.ServeHTTP(w, r)
	})
}

// authenticate middleware retrieves the token, looks up the session, and injects user info into the context.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, models.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		session, err := tokens.GetSession(app.redis, token, tokens.ScopeAuthentication)
		if err != nil {
			switch {
			case errors.Is(err, tokens.ErrSessionNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		//Fetching full user data from db now
		user, err := app.models.Users.GetByID(r.Context(), session.UserID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		// At this point, the session is valid. For now, we will create a basic User struct.
		// Later, we will fetch the full user details from the database.
		// user := &models.User{
		// 	ID:    session.UserID,
		// 	Roles: session.Roles,
		// }

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

// requireAuthenticatedUser middleware checks if a user is authenticated.
func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// requirePermission middleware checks if the authenticated user has the required permission.
func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		//Get the user's permissions from the database.
		permissions, err := app.models.Permissions.GetAllForUser(ctx, user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		//Check if the user has the required permission.
		_, ok := permissions[code]
		if !ok {
			app.notPermittedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}

	// Wrap with requireAuthenticatedUser to ensure the user is logged in first.
	return app.requireAuthenticatedUser(fn)
}

// context helpers for getting/setting user in the request context.
func (app *application) contextSetUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		// This should not happen if middleware is set up correctly.
		panic("missing user value in request context")
	}
	return user
}
