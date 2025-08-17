// File: backend/cmd/api/middleware.go

package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/cliffdoyle/internal/models"
	"github.com/cliffdoyle/internal/tokens"
)

// A custom contextKey type to avoid key collisions in context.
type contextKey string

const userContextKey = contextKey("user")

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

		// At this point, the session is valid. For now, we will create a basic User struct.
		// Later, we will fetch the full user details from the database.
		user := &models.User{
			ID:    session.UserID,
			Roles: session.Roles,
		}

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
