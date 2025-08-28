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

var (
	ErrDuplicateEmail = errors.New("email already exists")
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User, hotelID uuid.UUID) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	// Future methods: GetByID, Update, Delete
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User, hotelID uuid.UUID) error {
	query := `
        INSERT INTO users (hotel_id, email, password_hash, first_name, last_name)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at, version`

	args := []any{hotelID, user.Email, user.PasswordHash, user.FirstName, user.LastName}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Version)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        SELECT u.id, u.hotel_id, u.email, u.password_hash, u.first_name, u.last_name, u.is_active, u.created_at, u.updated_at, u.version,
               COALESCE(ARRAY_AGG(ro.name) FILTER (WHERE ro.name IS NOT NULL), '{}') as roles
        FROM users u
        LEFT JOIN users_roles ur ON u.id = ur.user_id
        LEFT JOIN roles ro ON ur.role_id = ro.id
        WHERE u.email = $1
        GROUP BY u.id`

	var user models.User
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.HotelID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.IsActive, &user.CreatedAt,
		&user.UpdatedAt, &user.Version, &user.Roles,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Add the implementation for GetByID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
        SELECT u.id, u.hotel_id, u.email, u.password_hash, u.first_name, u.last_name, u.is_active, u.created_at, u.updated_at, u.version,
               COALESCE(ARRAY_AGG(ro.name) FILTER (WHERE ro.name IS NOT NULL), '{}') as roles
        FROM users u
        LEFT JOIN users_roles ur ON u.id = ur.user_id
        LEFT JOIN roles ro ON ur.role_id = ro.id
        WHERE u.id = $1
        GROUP BY u.id`

	var user models.User
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.HotelID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.IsActive, &user.CreatedAt,
		&user.UpdatedAt, &user.Version, &user.Roles,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}
