package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/cliffdoyle/internal/repository"
	"github.com/cliffdoyle/internal/service"
	validator "github.com/cliffdoyle/internal/validation"
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

func (app *application) getReservationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	user := app.contextGetUser(r)

	res, err := app.services.Reservations.GetReservationByID(r.Context(), user.HotelID, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"reservation": res}, nil)
}

func (app *application) listReservationsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	qs := r.URL.Query()
	v := validator.New()

	var filters repository.ListReservationsFilters

	startDateStr := app.readString(qs, "startDate", "")
	endDateStr := app.readString(qs, "endDate", "")

	// Basic parsing. A production system would have more robust date validation.
	if startDateStr != "" {
		filters.StartDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			v.AddError("startDate", "must be in YYYY-MM-DD format")
		}
	}
	if endDateStr != "" {
		filters.EndDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			v.AddError("endDate", "must be in YYYY-MM-DD format")
		}
	}

	filters.Page = app.readInt(qs, "page", 1, v)
	filters.PageSize = app.readInt(qs, "pageSize", 50, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	reservations, err := app.services.Reservations.ListReservations(r.Context(), user.HotelID, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"reservations": reservations}, nil)
}

func (app *application) updateReservationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r); return
	}
	
	user := app.contextGetUser(r)

	var input service.ReservationUpdateDTO
	if err = app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err); return
	}

	res, err := app.services.Reservations.UpdateReservation(r.Context(), user.HotelID, id, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotAvailable):
			app.errorResponse(w, r, http.StatusConflict, err.Error())
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	
	app.writeJSON(w, http.StatusOK, envelope{"reservation": res}, nil)
}

func (app *application) cancelReservationHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	user := app.contextGetUser(r)

	err = app.services.Reservations.CancelReservation(r.Context(), user.HotelID, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "reservation cancelled successfully"}, nil)
}
