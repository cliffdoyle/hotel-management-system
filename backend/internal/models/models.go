// File: backend/internal/models/models.go

package models

import (
	"errors"
	"time"

	"github.com/cliffdoyle/internal/passwords"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user account in the system.
// @Description User account information along with their hotel and roles.
type User struct {
	ID           uuid.UUID `json:"id"`
	HotelID      uuid.UUID `json:"hotel_id"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"-"` // Omit from JSON responses
	FirstName    *string   `json:"first_name,omitempty"`
	LastName     *string   `json:"last_name,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Version      int       `json:"-"` // Used for optimistic locking
	Roles        []string  `json:"roles,omitempty"` // For RBAC
}

// Hotel represents a single hotel tenant in our multi-tenant system.
// @Description Represents a hotel property with its basic details.
type Hotel struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Role represents a user role (e.g., admin, manager).
// @Description A role that can be assigned to a user, defining their access level.
type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Permission represents a specific action a role can have.
// @Description A specific permission code, like 'reservations:read'.
type Permission struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
}

// AnonymousUser represents a user who is not authenticated.
var AnonymousUser = &User{}

// IsAnonymous checks if a User instance is the anonymous user.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// SetPassword hashes the plaintext password and stores it in the user struct.
func (u *User) SetPassword(plaintextPassword string) error {
	hash, err := passwords.Hash(plaintextPassword)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	return nil
}

// MatchesPassword checks if the plaintext password matches the stored hash.
func (u *User) MatchesPassword(plaintextPassword string) (bool, error) {
	match, err := passwords.Matches(plaintextPassword, u.PasswordHash)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil // Not an application error
		}
		return false, err
	}
	return match, nil
}
