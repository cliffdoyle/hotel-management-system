package repository

import (
	"context"
	"errors"

	"time"

	"github.com/cliffdoyle/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ListReservationsFilters holds the filtering parameters for listing reservations.
type ListReservationsFilters struct {
	StartDate time.Time
	EndDate   time.Time
	Page      int
	PageSize  int
}

type ReservationRepository interface {
	// Transactional methods
	CheckAvailability(ctx context.Context, tx pgx.Tx, hotelID, roomTypeID uuid.UUID, startDate, endDate time.Time) (bool, error)
	DecrementInventory(ctx context.Context, tx pgx.Tx, hotelID, roomTypeID uuid.UUID, startDate, endDate time.Time) error
	Create(ctx context.Context, tx pgx.Tx, res *models.Reservation) error

	// Non-transactional methods
	GetByID(ctx context.Context, hotelID, resID uuid.UUID) (*models.Reservation, error)
	List(ctx context.Context, hotelID uuid.UUID, filters ListReservationsFilters) ([]*models.Reservation, error)

	//Method for the availability calendar
	GetAvailabilityForDateRange(ctx context.Context, hotelID, roomTypeID uuid.UUID, startDate, endDate time.Time) (map[string]int, error)
	IncrementInventory(ctx context.Context, tx pgx.Tx, hotelID, roomTypeID uuid.UUID, startDate, endDate time.Time) error

	Update(ctx context.Context, tx pgx.Tx, res *models.Reservation) error
}

type reservationRepository struct {
	db *pgxpool.Pool
}

func NewReservationRepository(db *pgxpool.Pool) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(ctx context.Context, tx pgx.Tx, res *models.Reservation) error {
	query := `
		INSERT INTO reservations (hotel_id, guest_id, room_type_id, rate_plan_id, start_date, end_date, num_adults, num_children, status, total_cost_cents, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at, version`

	args := []any{res.HotelID, res.GuestID, res.RoomTypeID, res.RatePlanID, res.StartDate, res.EndDate, res.NumAdults, res.NumChildren, res.Status, res.TotalCostCents, res.Notes}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return tx.QueryRow(ctx, query, args...).Scan(&res.ID, &res.CreatedAt, &res.UpdatedAt, &res.Version)
}

func (r *reservationRepository) GetByID(ctx context.Context, hotelID, resID uuid.UUID) (*models.Reservation, error) {
	query := `
		SELECT 
			id, hotel_id, guest_id, room_type_id, rate_plan_id,
			start_date, end_date, num_adults, num_children, status,
			total_cost_cents, notes, created_at, updated_at, version
		FROM reservations
		WHERE id = $1 AND hotel_id = $2`

	var res models.Reservation
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, query, resID, hotelID).Scan(
		&res.ID, &res.HotelID, &res.GuestID, &res.RoomTypeID, &res.RatePlanID,
		&res.StartDate, &res.EndDate, &res.NumAdults, &res.NumChildren, &res.Status,
		&res.TotalCostCents, &res.Notes, &res.CreatedAt, &res.UpdatedAt, &res.Version,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &res, nil
}

func (r *reservationRepository) CheckAvailability(ctx context.Context, tx pgx.Tx, hotelID, roomTypeID uuid.UUID, startDate, endDate time.Time) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM inventory_levels
		WHERE hotel_id = $1
		  AND room_type_id = $2
		  AND date >= $3 AND date < $4
		  AND available_rooms > 0`

	var count int
	err := tx.QueryRow(ctx, query, hotelID, roomTypeID, startDate, endDate).Scan(&count)
	if err != nil {
		return false, err
	}

	numberOfNights := int(endDate.Sub(startDate).Hours() / 24)
	return count == numberOfNights, nil
}

func (r *reservationRepository) DecrementInventory(ctx context.Context, tx pgx.Tx, hotelID, roomTypeID uuid.UUID, startDate, endDate time.Time) error {
	query := `
		UPDATE inventory_levels
		SET available_rooms = available_rooms - 1
		WHERE hotel_id = $1
		  AND room_type_id = $2
		  AND date >= $3 AND date < $4
          AND available_rooms > 0` // Added safety check

	tag, err := tx.Exec(ctx, query, hotelID, roomTypeID, startDate, endDate)
	if err != nil {
		return err
	}

	// Safety check: ensure the number of updated rows matches the number of nights.
	// If it doesn't, it means another transaction took the last room in the middle of our transaction (a race condition).
	numberOfNights := int(endDate.Sub(startDate).Hours() / 24)
	if tag.RowsAffected() != int64(numberOfNights) {
		return errors.New("failed to decrement inventory for all nights, potential race condition")
	}

	return nil
}

func (r *reservationRepository) IncrementInventory(ctx context.Context, tx pgx.Tx, hotelID, roomTypeID uuid.UUID, startDate, endDate time.Time) error {
	query := `
		UPDATE inventory_levels
		SET available_rooms = available_rooms + 1
		WHERE hotel_id = $1
		  AND room_type_id = $2
		  AND date >= $3 AND date < $4`

	_, err := tx.Exec(ctx, query, hotelID, roomTypeID, startDate, endDate)
	return err
}

func (r *reservationRepository) GetAvailabilityForDateRange(ctx context.Context, hotelID, roomTypeID uuid.UUID, startDate, endDate time.Time) (map[string]int, error) {
	query := `
		SELECT date, available_rooms
		FROM inventory_levels
		WHERE hotel_id = $1
		  AND room_type_id = $2
		  AND date >= $3 AND date < $4`

	rows, err := r.db.Query(ctx, query, hotelID, roomTypeID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	availability := make(map[string]int)
	for rows.Next() {
		var date time.Time
		var availableRooms int
		if err := rows.Scan(&date, &availableRooms); err != nil {
			return nil, err
		}
		// Format the date as YYYY-MM-DD for a consistent key in our map
		dateStr := date.Format("2006-01-02")
		availability[dateStr] = availableRooms
	}

	return availability, rows.Err()
}

func (r *reservationRepository) Update(ctx context.Context, tx pgx.Tx, res *models.Reservation) error {
	query := `
		UPDATE reservations
		SET 
			start_date = $1,
			end_date = $2,
			num_adults = $3,
			total_cost_cents = $4, -- assuming we've recalculated it
			notes = $5,
			updated_at = NOW(),
			version = version + 1
		WHERE id = $6 AND version = $7`

	args := []any{
		res.StartDate, res.EndDate, res.NumAdults,
		res.TotalCostCents, res.Notes,
		res.ID, res.Version,
	}

	cmdTag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	// Optimistic locking check
	if cmdTag.RowsAffected() == 0 {
		return errors.New("edit conflict: reservation has been modified by another user")
	}

	return nil
}

func (r *reservationRepository) List(ctx context.Context, hotelID uuid.UUID, filters ListReservationsFilters) ([]*models.Reservation, error) {
	// This query finds all reservations that *overlap* with the given date range.
	// (start_date, end_date) OVERLAPS (filters.StartDate, filters.EndDate)
	query := `
		SELECT
			r.id, r.guest_id, r.room_type_id, r.rate_plan_id, r.start_date, r.end_date, r.status,
			g.id, g.first_name, g.last_name, g.email
		FROM reservations r
		INNER JOIN guests g ON r.guest_id = g.id
		WHERE r.hotel_id = $1
		AND (r.start_date, r.end_date) OVERLAPS ($2, $3)
		ORDER BY r.start_date
		LIMIT $4 OFFSET $5`

	offset := (filters.Page - 1) * filters.PageSize

	rows, err := r.db.Query(ctx, query, hotelID, filters.StartDate, filters.EndDate, filters.PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*models.Reservation
	for rows.Next() {
		var res models.Reservation
		var guest models.Guest
		err := rows.Scan(
			&res.ID, &res.GuestID, &res.RoomTypeID, &res.RatePlanID, &res.StartDate, &res.EndDate, &res.Status,
			&guest.ID, &guest.FirstName, &guest.LastName, &guest.Email,
		)
		if err != nil {
			return nil, err
		}
		res.Guest = &guest // Embed guest details for the UI
		reservations = append(reservations, &res)
	}

	return reservations, rows.Err()
}
