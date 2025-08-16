// File: backend/cmd/api/errors.go

package main

import (
	"fmt"
	"net/http"
)

// logError is a generic helper for logging an error message.
func (app *application) logError(r *http.Request, err error) {
	app.logger.Error(err.Error(), "request_method", r.Method, "request_url", r.URL.String())
}

// errorResponse is a generic helper for sending JSON-formatted error messages to the client
// with a given status code.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// serverErrorResponse will be used when our server encounters an unexpected problem.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse will be used to send a 404 Not Found status code and JSON response.
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse will be used to send a 405 Method Not Allowed status code and JSON response.
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// badRequestResponse will be used for malformed requests from the client.
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// failedValidationResponse will be used for validation errors on DTOs.
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "you do not have the necessary permissions to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}
