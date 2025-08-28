package main

import (
	"errors"
	"net/http"

	"github.com/cliffdoyle/internal/repository"
	"github.com/cliffdoyle/internal/service"
	validator "github.com/cliffdoyle/internal/validation"
)

func (app *application) createGuestHandler(w http.ResponseWriter, r *http.Request) {
	var input service.GuestCreateDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	guest, err := app.services.Guests.CreateGuest(r.Context(), input, user.HotelID)
	if err != nil {
		var validationErr *validator.ValidationError
		switch {
		case errors.As(err, &validationErr):
			app.failedValidationResponse(w, r, validationErr.Errors)
		case errors.Is(err, repository.ErrDuplicateEmail):
			app.failedValidationResponse(w, r, map[string]string{"email": "a guest with this email already exists"})
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"guest": guest}, nil)
}

func (app *application) getGuestHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)

	guest, err := app.services.Guests.GetGuestByID(r.Context(), user.HotelID, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"guest": guest}, nil)
}

func (app *application) listGuestsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	var filters repository.ListGuestsFilters
	v := validator.New()
	qs := r.URL.Query()

	filters.Name = app.readString(qs, "name", "")
	filters.Email = app.readString(qs, "email", "")
	filters.Page = app.readInt(qs, "page", 1, v)
	filters.PageSize = app.readInt(qs, "pageSize", 20, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	guests, err := app.services.Guests.ListGuests(r.Context(), user.HotelID, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"guests": guests}, nil)
}
