package service

import (
	"context"
	"errors"
	"time"

	"github.com/cliffdoyle/internal/models"
	"github.com/cliffdoyle/internal/repository"
	validator "github.com/cliffdoyle/internal/validation"
	"github.com/google/uuid"
)

var (
	ErrPriceNotAvailable = errors.New("a price is not available for one or more nights in the requested date range")
)

type RatesCreateDTO struct {
	RoomTypeID uuid.UUID `json:"room_type_id"`
	RatePlanID uuid.UUID `json:"rate_plan_id"`
	StartDate  string    `json:"start_date"` // YYYY-MM-DD
	EndDate    string    `json:"end_date"`   // YYYY-MM-DD
	PriceCents int       `json:"price_cents"`
}

type PriceQuoteResponse struct {
	TotalCents     int            `json:"total_cents"`
	NightlyRates   []*models.Rate `json:"nightly_rates"`
	Currency       string         `json:"currency"`
	NumberOfNights int            `json:"number_of_nights"`
}

type RateService interface {
	CreateRatesForDateRange(ctx context.Context, dto RatesCreateDTO, hotelID uuid.UUID) (int64, error)
	GetPriceQuote(ctx context.Context, hotelID, roomTypeID, ratePlanID uuid.UUID, startDateStr, endDateStr string) (*PriceQuoteResponse, error)
}

type rateService struct {
	rateRepo repository.RateRepository
}

func NewRateService(rateRepo repository.RateRepository) RateService {
	return &rateService{rateRepo: rateRepo}
}

func (s *rateService) CreateRatesForDateRange(ctx context.Context, dto RatesCreateDTO, hotelID uuid.UUID) (int64, error) {
	startDate, err := time.Parse("2006-01-02", dto.StartDate)
	if err != nil {
		return 0, errors.New("invalid start_date format, use YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", dto.EndDate)
	if err != nil {
		return 0, errors.New("invalid end_date format, use YYYY-MM-DD")
	}

	v := validator.New()
	v.Check(endDate.After(startDate), "end_date", "must be after start_date")
	v.Check(dto.PriceCents >= 0, "price_cents", "must be a positive value")
	if !v.Valid() {
		return 0, &validator.ValidationError{Errors: v.Errors}
	}

	var ratesToInsert []*models.Rate
	// Loop through each day in the range and create a Rate model
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		rate := &models.Rate{
			HotelID:    hotelID,
			RoomTypeID: dto.RoomTypeID,
			RatePlanID: dto.RatePlanID,
			Date:       d,
			PriceCents: dto.PriceCents,
		}
		ratesToInsert = append(ratesToInsert, rate)
	}

	if len(ratesToInsert) == 0 {
		return 0, errors.New("no rates to insert for the given date range")
	}

	return s.rateRepo.BulkInsert(ctx, ratesToInsert)
}

func (s *rateService) GetPriceQuote(ctx context.Context, hotelID, roomTypeID, ratePlanID uuid.UUID, startDateStr, endDateStr string) (*PriceQuoteResponse, error) {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, errors.New("invalid start_date format, use YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, errors.New("invalid end_date format, use YYYY-MM-DD")
	}

	v := validator.New()
	v.Check(endDate.After(startDate), "end_date", "must be after start_date")
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	numberOfNights := int(endDate.Sub(startDate).Hours() / 24)

	nightlyRates, err := s.rateRepo.FindForStay(ctx, hotelID, roomTypeID, ratePlanID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Business rule: if we don't have a price for every single night, the stay is not bookable.
	if len(nightlyRates) != numberOfNights {
		return nil, ErrPriceNotAvailable
	}

	totalCents := 0
	for _, rate := range nightlyRates {
		totalCents += rate.PriceCents
	}

	return &PriceQuoteResponse{
		TotalCents:     totalCents,
		NightlyRates:   nightlyRates,
		Currency:       "USD", // Hardcoded for now, should come from hotel settings
		NumberOfNights: numberOfNights,
	}, nil
}
