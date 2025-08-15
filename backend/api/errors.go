package main

import (
	"fmt"
	"log"
	"net/http"
)

// Custom error types for different scenarios

// AppError represents application-specific errors
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Predefined error instances

var (
	// Authentication errors
	ErrInvalidCredentials = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid credentials",
	}
	ErrUnauthorized = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Unauthorized access",
	}
	ErrTokenExpired = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Token has expired",
	}
	ErrInvalidToken = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "Invalid token",
	}

	// Validation errors
	ErrValidationFailed = &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: "Validation failed",
	}
	ErrInvalidInput = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid input data",
	}

	// Resource errors
	ErrNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "Resource not found",
	}
	ErrConflict = &AppError{
		Code:    http.StatusConflict,
		Message: "Resource already exists",
	}

	// Server errors
	ErrInternalServer = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	}
	ErrDatabaseConnection = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Database connection error",
	}
	ErrRedisConnection = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Redis connection error",
	}

	// Rate limiting errors
	ErrRateLimitExceeded = &AppError{
		Code:    http.StatusTooManyRequests,
		Message: "Rate limit exceeded",
	}

	// Permission errors
	ErrForbidden = &AppError{
		Code:    http.StatusForbidden,
		Message: "Access forbidden",
	}
	ErrInsufficientPermissions = &AppError{
		Code:    http.StatusForbidden,
		Message: "Insufficient permissions",
	}
)

// Error handling methods for the Application struct

// logError logs error details for debugging
func (app *Application) logError(r *http.Request, err error) {
	log.Printf("ERROR: %s %s - %v", r.Method, r.URL.Path, err)
}

// errorResponse sends error response to client
func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"success": false, "error": message}
	
	err := app.writeJSON(w, status, env)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// serverErrorResponse handles 500 internal server errors
func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	app.errorResponse(w, r, http.StatusInternalServerError, "The server encountered a problem and could not process your request")
}

// notFoundResponse handles 404 not found errors
func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusNotFound, "The requested resource could not be found")
}

// methodNotAllowedResponse handles 405 method not allowed errors
func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusMethodNotAllowed, fmt.Sprintf("The %s method is not supported for this resource", r.Method))
}

// badRequestResponse handles 400 bad request errors
func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// validationErrorResponse handles validation errors
func (app *Application) validationErrorResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// unauthorizedResponse handles 401 unauthorized errors
func (app *Application) unauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusUnauthorized, "You must be authenticated to access this resource")
}

// forbiddenResponse handles 403 forbidden errors
func (app *Application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusForbidden, "Your account doesn't have permission to access this resource")
}

// rateLimitExceededResponse handles 429 rate limit exceeded errors
func (app *Application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusTooManyRequests, "Rate limit exceeded")
}

// handleAppError handles custom AppError types
func (app *Application) handleAppError(w http.ResponseWriter, r *http.Request, appErr *AppError) {
	if appErr.Err != nil {
		app.logError(r, appErr.Err)
	}
	app.errorResponse(w, r, appErr.Code, appErr.Message)
}
