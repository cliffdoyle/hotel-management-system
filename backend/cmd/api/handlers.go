// File: backend/cmd/api/handlers.go

package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cliffdoyle/internal/tokens"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     "1.0.0",
		},
		"dependencies": map[string]string{
			"database": "OK",
			"redis":    "OK",
		},
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Use app.db directly to ping the database.
	if err := app.db.Ping(ctx); err != nil {
		data["dependencies"].(map[string]string)["database"] = "Error: " + err.Error()
		data["status"] = "degraded"
	}

	if err := app.redis.Ping(ctx).Err(); err != nil {
		data["dependencies"].(map[string]string)["redis"] = "Error: " + err.Error()
		data["status"] = "degraded"
	}

	statusCode := http.StatusOK
	if data["status"] != "available" {
		statusCode = http.StatusServiceUnavailable
	}

	err := app.writeJSON(w, statusCode, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// refreshTokenHandler validates a refresh token and issues a new pair of tokens.
func (app *application) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// For now, validation on the DTO is simple.
	// In the next issue, we'll add a proper validator.

	// Look up the session for the refresh token.
	session, err := tokens.GetSession(app.redis, input.RefreshToken, tokens.ScopeRefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, tokens.ErrSessionNotFound):
			app.invalidAuthenticationTokenResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// The refresh token is valid. Immediately delete the old one.
	err = tokens.DeleteSession(app.redis, input.RefreshToken, tokens.ScopeRefreshToken)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Issue a new pair of tokens.
	// NOTE: For now, roles are an empty slice. We will fetch real roles during login in Issue 5.
	newAuthToken, err := tokens.GenerateToken(app.redis, session.UserID, 30*time.Minute, tokens.ScopeAuthentication, []string{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	newRefreshToken, err := tokens.GenerateToken(app.redis, session.UserID, 24*time.Hour*30, tokens.ScopeRefreshToken, []string{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the new tokens to the client.
	response := envelope{
		"authentication_token": newAuthToken,
		"refresh_token":        newRefreshToken.Plaintext,
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// metricsHandler creates a handler that serves metrics from our CUSTOM registry.
func (app *application) metricsHandler() http.Handler {
	return promhttp.HandlerFor(app.metrics_reg, promhttp.HandlerOpts{})
}
