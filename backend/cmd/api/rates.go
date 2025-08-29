package main

import (
	"errors"
	"net/http"

	"github.com/cliffdoyle/internal/service"
	validator "github.com/cliffdoyle/internal/validation"
	"github.com/google/uuid"
)

func (app *application) createRatesHandler(w http.ResponseWriter, r *http.Request) {
	var input service.RatesCreateDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	count, err := app.services.Rates.CreateRatesForDateRange(r.Context(), input, user.HotelID)
	if err != nil {
		var validationErr *validator.ValidationError
		if errors.As(err, &validationErr) {
			app.failedValidationResponse(w, r, validationErr.Errors)
		} else {
			app.badRequestResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"message": "rates created successfully", "records_created": count}, nil)
}

func (app *application) getPriceQuoteHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	v := validator.New()

	roomTypeIDStr := app.readString(qs, "roomTypeId", "")
	ratePlanIDStr := app.readString(qs, "ratePlanId", "")
	startDateStr := app.readString(qs, "startDate", "")
	endDateStr := app.readString(qs, "endDate", "")

	v.Check(roomTypeIDStr != "", "roomTypeId", "must be provided")
	v.Check(ratePlanIDStr != "", "ratePlanId", "must be provided")

	roomTypeID, err := uuid.Parse(roomTypeIDStr)
	if err != nil {
		v.AddError("roomTypeId", "must be a valid UUID")
	}

	ratePlanID, err := uuid.Parse(ratePlanIDStr)
	if err != nil {
		v.AddError("ratePlanId", "must be a valid UUID")
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// For a quote, we need the hotel ID. This might come from the URL in a multi-tenant system
	// e.g., grand-hotel.pms.com. For now, we'll require authentication.
	user := app.contextGetUser(r)

	quote, err := app.services.Rates.GetPriceQuote(r.Context(), user.HotelID, roomTypeID, ratePlanID, startDateStr, endDateStr)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPriceNotAvailable):
			app.errorResponse(w, r, http.StatusNotFound, err.Error())
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"price_quote": quote}, nil)
}

// NOTE: RatePlan CRUD handlers would be created in a similar fashion.
