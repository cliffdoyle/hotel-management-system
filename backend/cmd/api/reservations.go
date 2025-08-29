package main

import (
	"errors"
	"net/http"

	"github.com/cliffdoyle/internal/service"
)

func (app *application) createReservationHandler(w http.ResponseWriter, r *http.Request) {
	var input service.ReservationCreateDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	res, err := app.services.Reservations.CreateReservation(r.Context(), input, user.HotelID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotAvailable):
			app.errorResponse(w, r, http.StatusConflict, err.Error()) // 409 Conflict
		case errors.Is(err, service.ErrPriceNotAvailable):
			app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"reservation": res}, nil)
}
