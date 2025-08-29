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

type ReservationRepository interface {
	// Transactional methods
	CheckAvailability(ctx context.Context, tx pgx.Tx, hotelID, roomTypeID uuid.UUID, startDate, endDate time.Time) (bool, error)
	DecrementInventory(ctx context.Context, tx pgx.Tx, hotelID, roomTypeID uuid.UUID, startDate, endDate time.Time) error
	Create(ctx context.Context, tx pgx.Tx, res *models.Reservation) error

	// Non-transactional methods
	GetByID(ctx context.Context, hotelID, resID uuid.UUID) (*models.Reservation, error)
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
