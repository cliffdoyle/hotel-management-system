package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// JSON Response helpers

// writeJSON writes JSON response with proper headers and error handling
func (app *Application) writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

// readJSON reads and parses JSON from request body with proper error handling
func (app *Application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Limit request body size to 1MB
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Create decoder and disallow unknown fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the JSON
	if err := dec.Decode(dst); err != nil {
		return err
	}

	// Check if there's additional data after the JSON object
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("body must only contain a single JSON object")
	}

	return nil
}

// Parameter extraction helpers

// readIDParam extracts and validates ID parameter from URL
func (app *Application) readIDParam(ps httprouter.Params) (int64, error) {
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("invalid id parameter")
	}
	return id, nil
}

// readStringParam extracts string parameter from URL
func (app *Application) readStringParam(ps httprouter.Params, key string) string {
	return ps.ByName(key)
}

// Query parameter helpers

// readString extracts string query parameter with default value
func (app *Application) readString(qs map[string][]string, key string, defaultValue string) string {
	s := qs[key]
	if len(s) == 0 {
		return defaultValue
	}
	return s[0]
}

// readInt extracts integer query parameter with default value
func (app *Application) readInt(qs map[string][]string, key string, defaultValue int) int {
	s := qs[key]
	if len(s) == 0 {
		return defaultValue
	}

	i, err := strconv.Atoi(s[0])
	if err != nil {
		return defaultValue
	}
	return i
}

// Response envelope for consistent API responses
type envelope map[string]interface{}

// writeSuccessResponse writes a successful response with data
func (app *Application) writeSuccessResponse(w http.ResponseWriter, status int, data interface{}) error {
	return app.writeJSON(w, status, envelope{"success": true, "data": data})
}

// writeErrorResponse writes an error response
func (app *Application) writeErrorResponse(w http.ResponseWriter, status int, message string) error {
	return app.writeJSON(w, status, envelope{"success": false, "error": message})
}

// writeValidationErrorResponse writes validation error response
func (app *Application) writeValidationErrorResponse(w http.ResponseWriter, errors map[string]string) error {
	return app.writeJSON(w, http.StatusUnprocessableEntity, envelope{
		"success": false,
		"error":   "validation failed",
		"errors":  errors,
	})
}
