// File: backend/internal/models/models.go

package models

import (
	"time"

	"github.com/google/uuid"
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