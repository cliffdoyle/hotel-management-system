package main

import (
	"errors"
	"net/http"

	"github.com/cliffdoyle/internal/repository"
	"github.com/cliffdoyle/internal/service"
	validator "github.com/cliffdoyle/internal/validation"
	"github.com/google/uuid"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input service.UserRegisterDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// IMPORTANT: Replace with a real hotel UUID from your database.
	hotelID, _ := uuid.Parse("694b475a-31ca-4ec5-a092-223e71ff69e2")

	user, err := app.services.Users.Register(r.Context(), input, hotelID)
	if err != nil {
		var validationErr *validator.ValidationError 
		switch {
		case errors.Is(err, repository.ErrDuplicateEmail):
			app.failedValidationResponse(w, r, map[string]string{"email": "a user with this email address already exists"})
		case errors.As(err, &validationErr):
			app.failedValidationResponse(w, r, validationErr.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input service.UserLoginDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	response, err := app.services.Users.Login(r.Context(), input)
	if err != nil {
		var validationErr *validator.ValidationError
		switch {
		case errors.Is(err, repository.ErrInvalidCredentials):
			app.errorResponse(w, r, http.StatusUnauthorized, "invalid email or password")
		case errors.As(err, &validationErr):
			app.failedValidationResponse(w, r, validationErr.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"tokens": response}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
