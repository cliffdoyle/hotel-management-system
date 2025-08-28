package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/cliffdoyle/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ListGuestsFilters struct {
	Name     string
	Email    string
	Page     int
	PageSize int
}

type GuestRepository interface {
	Create(ctx context.Context, guest *models.Guest) error
	GetByID(ctx context.Context, hotelID, guestID uuid.UUID) (*models.Guest, error)
	List(ctx context.Context, hotelID uuid.UUID, filters ListGuestsFilters) ([]*models.Guest, error)
	// Update coming in next issue to show partial updates
}

type guestRepository struct {
	db *pgxpool.Pool
}

func NewGuestRepository(db *pgxpool.Pool) GuestRepository {
	return &guestRepository{db: db}
}

func (r *guestRepository) Create(ctx context.Context, guest *models.Guest) error {
	query := `
		INSERT INTO guests (hotel_id, first_name, last_name, email, phone, preferences)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at, version`

	args := []any{guest.HotelID, guest.FirstName, guest.LastName, guest.Email, guest.Phone, guest.Preferences}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, query, args...).Scan(&guest.ID, &guest.CreatedAt, &guest.UpdatedAt, &guest.Version)
	if err != nil {
		if strings.Contains(err.Error(), "guests_hotel_id_email_key") {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (r *guestRepository) GetByID(ctx context.Context, hotelID, guestID uuid.UUID) (*models.Guest, error) {
	query := `
		SELECT id, hotel_id, first_name, last_name, email, phone, 
		       loyalty_member_number, loyalty_tier, loyalty_points, preferences,
		       metadata, created_at, updated_at, version
		FROM guests
		WHERE id = $1 AND hotel_id = $2`

	var guest models.Guest
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, query, guestID, hotelID).Scan(
		&guest.ID, &guest.HotelID, &guest.FirstName, &guest.LastName, &guest.Email, &guest.Phone,
		&guest.LoyaltyMemberNumber, &guest.LoyaltyTier, &guest.LoyaltyPoints, &guest.Preferences,
		&guest.Metadata, &guest.CreatedAt, &guest.UpdatedAt, &guest.Version,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &guest, nil
}

func (r *guestRepository) List(ctx context.Context, hotelID uuid.UUID, filters ListGuestsFilters) ([]*models.Guest, error) {
	// ILIKE for case-insensitive search
	query := `
		SELECT id, first_name, last_name, email, phone
		FROM guests
		WHERE hotel_id = $1
		  AND (LOWER(first_name || ' ' || last_name) LIKE LOWER($2) OR $2 = '')
		  AND (LOWER(email) = LOWER($3) OR $3 = '')
		ORDER BY last_name ASC, first_name ASC
		LIMIT $4 OFFSET $5`

	nameFilter := ""
	if filters.Name != "" {
		nameFilter = "%" + filters.Name + "%"
	}

	emailFilter := ""
	if filters.Email != "" {
		emailFilter = filters.Email
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	offset := (filters.Page - 1) * filters.PageSize
	rows, err := r.db.Query(ctx, query, hotelID, nameFilter, emailFilter, filters.PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guests []*models.Guest
	for rows.Next() {
		var guest models.Guest
		err := rows.Scan(&guest.ID, &guest.FirstName, &guest.LastName, &guest.Email, &guest.Phone)
		if err != nil {
			return nil, err
		}
		guests = append(guests, &guest)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return guests, nil
}
