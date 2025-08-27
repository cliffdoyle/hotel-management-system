package main

import (
	"errors"
	"net/http"

	"github.com/cliffdoyle/internal/repository"
	"github.com/cliffdoyle/internal/service"
	validator "github.com/cliffdoyle/internal/validation"
)

func (app *application) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	var input service.RoomCreateDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// For now, we get the user from the context to determine the hotel.
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	room, err := app.services.Rooms.CreateRoom(r.Context(), input, user.HotelID)
	if err != nil {
		var validationErr *validator.ValidationError
		switch {
		case errors.As(err, &validationErr):
			app.failedValidationResponse(w, r, validationErr.Errors)
		case errors.Is(err, repository.ErrDuplicateRoomNumber):
			app.failedValidationResponse(w, r, map[string]string{"room_number": "room number already exists"})
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"room": room}, nil)
}

func (app *application) listRoomsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	rooms, err := app.services.Rooms.GetRoomsForHotel(r.Context(), user.HotelID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"rooms": rooms}, nil)
}
