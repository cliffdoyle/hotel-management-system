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

type ReservationUpdateDTO struct {
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	NumAdults   int    `json:"num_adults"`
	// Other fields like notes, etc. can be added here
}

type ReservationService interface {
	CreateReservation(ctx context.Context, dto ReservationCreateDTO, hotelID uuid.UUID) (*models.Reservation, error)
	GetReservationByID(ctx context.Context, hotelID, resID uuid.UUID) (*models.Reservation, error)
	ListReservations(ctx context.Context, hotelID uuid.UUID, filters repository.ListReservationsFilters) ([]*models.Reservation, error) 
	// REVISED Update method
	UpdateReservation(ctx context.Context, hotelID, resID uuid.UUID, dto ReservationUpdateDTO) (*models.Reservation, error)
    // NEW method for the calendar
	CheckAvailability(ctx context.Context, hotelID, roomTypeID uuid.UUID, startDateStr, endDateStr string) (map[string]int, error)
	CancelReservation(ctx context.Context, hotelID, resID uuid.UUID) error
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

func (s *reservationService) GetReservationByID(ctx context.Context, hotelID, resID uuid.UUID) (*models.Reservation, error) {
	return s.resRepo.GetByID(ctx, hotelID, resID)
}


func (s *reservationService) ListReservations(ctx context.Context, hotelID uuid.UUID, filters repository.ListReservationsFilters) ([]*models.Reservation, error) {
	// Apply default pagination and date ranges if empty
	if filters.StartDate.IsZero() {
		filters.StartDate = time.Now()
	}
	if filters.EndDate.IsZero() {
		filters.EndDate = filters.StartDate.AddDate(0, 0, 30) // Default to a 30-day window
	}
	if filters.Page <= 0 { filters.Page = 1 }
	if filters.PageSize <= 0 { filters.PageSize = 50 }
	
	return s.resRepo.List(ctx, hotelID, filters)
}


func (s *reservationService) CheckAvailability(ctx context.Context, hotelID, roomTypeID uuid.UUID, startDateStr, endDateStr string) (map[string]int, error) {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil { return nil, errors.New("invalid start_date format") }
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil { return nil, errors.New("invalid end_date format") }

	// Validation logic could go here
	
	return s.resRepo.GetAvailabilityForDateRange(ctx, hotelID, roomTypeID, startDate, endDate)
}


func (s *reservationService) UpdateReservation(ctx context.Context, hotelID, resID uuid.UUID, dto ReservationUpdateDTO) (*models.Reservation, error) {
	// Validation of DTO...
	newStartDate, _ := time.Parse("2006-01-02", dto.StartDate)
	newEndDate, _   := time.Parse("2006-01-02", dto.EndDate)

	// A reservation update is a highly sensitive atomic transaction.
	tx, err := s.db.Begin(ctx)
	if err != nil { return nil, err }
	defer tx.Rollback(ctx)

	// 1. Get the original, current state of the reservation
	originalRes, err := s.resRepo.GetByID(ctx, hotelID, resID)
	if err != nil { return nil, err } // Handles not found

	// 2. Return the inventory for the OLD dates
	err = s.resRepo.IncrementInventory(ctx, tx, hotelID, originalRes.RoomTypeID, originalRes.StartDate, originalRes.EndDate)
	if err != nil { return nil, err }
	
	// 3. Check availability for the NEW dates
	available, err := s.resRepo.CheckAvailability(ctx, tx, hotelID, originalRes.RoomTypeID, newStartDate, newEndDate)
	if err != nil { return nil, err }
	if !available { return nil, ErrNotAvailable }
	
	// 4. Claim the inventory for the NEW dates
	err = s.resRepo.DecrementInventory(ctx, tx, hotelID, originalRes.RoomTypeID, newStartDate, newEndDate)
	if err != nil { return nil, err }
	
	// 5. TODO: Recalculate the price. For now, we will keep it the same.
	// This would involve calling the rateSvc.GetPriceQuote for the new dates.
	
	// 6. Update the reservation record in the database
    // For this, we'll need an Update method in the repository. Let's add it.
    originalRes.StartDate = newStartDate
    originalRes.EndDate = newEndDate
    originalRes.NumAdults = dto.NumAdults
	err = s.resRepo.Update(ctx, tx, originalRes)
    if err != nil { return nil, err }

	// 7. If everything succeeded, commit the transaction
	if err = tx.Commit(ctx); err != nil { return nil, err }

	return originalRes, nil
}


func (s *reservationService) CancelReservation(ctx context.Context, hotelID, resID uuid.UUID) error {
	// Cancellation is a transactional operation: update status AND return inventory.
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	
	// 1. Get the original reservation to know which rooms/dates to return to inventory.
	// We need the full reservation object.
	res, err := s.resRepo.GetByID(ctx, hotelID, resID)
	if err != nil {
		return err // Handles ErrRecordNotFound
	}

	// Business rule: Cannot cancel a reservation that's already in a final state.
	if res.Status == models.StatusCheckedOut || res.Status == models.StatusCancelled {
		return errors.New("cannot cancel a reservation with status " + string(res.Status))
	}
	
	// 2. Return the inventory for the dates of the stay.
	// This operation now happens inside the transaction.
	err = s.resRepo.IncrementInventory(ctx, tx, res.HotelID, res.RoomTypeID, res.StartDate, res.EndDate)
	if err != nil {
		return err
	}
	
    // --- START OF THE FIX ---
	
	// 3. Update the reservation's status in the Go struct.
	res.Status = models.StatusCancelled
	
	// 4. Pass the entire updated reservation object to the Update method.
    // The repository's Update method runs within the same transaction.
	err = s.resRepo.Update(ctx, tx, res)
	if err != nil {
		return err
	}

    // --- END OF THE FIX ---
	
	// 5. If everything succeeds, commit the transaction.
	return tx.Commit(ctx)
}