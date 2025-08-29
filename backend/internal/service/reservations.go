package service

import (
	"context"
	"errors"
	"time"

	"github.com/cliffdoyle/internal/models"
	"github.com/cliffdoyle/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotAvailable = errors.New("the selected room type is not available for the given dates")
)

type ReservationCreateDTO struct {
	GuestID    uuid.UUID `json:"guest_id"`
	RoomTypeID uuid.UUID `json:"room_type_id"`
	RatePlanID uuid.UUID `json:"rate_plan_id"`
	StartDate  string    `json:"start_date"`
	EndDate    string    `json:"end_date"`
	NumAdults  int       `json:"num_adults"`
	// ... other fields
}

type ReservationService interface {
	CreateReservation(ctx context.Context, dto ReservationCreateDTO, hotelID uuid.UUID) (*models.Reservation, error)
}

type reservationService struct {
	resRepo repository.ReservationRepository
	rateSvc RateService
	db      *pgxpool.Pool
}

func NewReservationService(resRepo repository.ReservationRepository, rateSvc RateService, db *pgxpool.Pool) ReservationService {
	return &reservationService{resRepo: resRepo, rateSvc: rateSvc, db: db}
}

func (s *reservationService) CreateReservation(ctx context.Context, dto ReservationCreateDTO, hotelID uuid.UUID) (*models.Reservation, error) {
	// Parse dates from DTO...
	// Validate DTO...

	startDate, _ := time.Parse("2006-01-02", dto.StartDate)
	endDate, _ := time.Parse("2006-01-02", dto.EndDate)

	// --- ATOMIC TRANSACTION ---
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // Safety net

	// 1. Get the price quote within the transaction
	quote, err := s.rateSvc.GetPriceQuote(ctx, hotelID, dto.RoomTypeID, dto.RatePlanID, dto.StartDate, dto.EndDate)
	if err != nil {
		return nil, err // Price not available
	}

	// 2. Check for room availability within the transaction
	available, err := s.resRepo.CheckAvailability(ctx, tx, hotelID, dto.RoomTypeID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, ErrNotAvailable
	}

	// 3. Decrement the inventory count for each night of the stay
	err = s.resRepo.DecrementInventory(ctx, tx, hotelID, dto.RoomTypeID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 4. Create the reservation record
	res := &models.Reservation{
		HotelID:        hotelID,
		GuestID:        dto.GuestID,
		RoomTypeID:     dto.RoomTypeID,
		RatePlanID:     dto.RatePlanID,
		StartDate:      startDate,
		EndDate:        endDate,
		NumAdults:      dto.NumAdults,
		Status:         models.StatusConfirmed,
		TotalCostCents: quote.TotalCents,
	}

	err = s.resRepo.Create(ctx, tx, res)
	if err != nil {
		return nil, err
	}

	// 5. If all steps succeed, commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return res, nil
}
