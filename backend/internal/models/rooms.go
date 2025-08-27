package models

import (
	"time"

	"github.com/google/uuid"
)

// RoomStatus represents the possible states of a room.
type RoomStatus string

const (
	StatusAvailableClean   RoomStatus = "AVAILABLE_CLEAN"
	StatusAvailableDirty   RoomStatus = "AVAILABLE_DIRTY"
	StatusOccupied         RoomStatus = "OCCUPIED"
	StatusOutOfService     RoomStatus = "OUT_OF_SERVICE"
)

// RoomType defines a category of room (e.g., "King Suite", "Standard Double").
type RoomType struct {
	ID          uuid.UUID `json:"id"`
	HotelID     uuid.UUID `json:"-"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Capacity    int       `json:"capacity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Room represents an individual, physical room in the hotel.
type Room struct {
	ID          uuid.UUID  `json:"id"`
	HotelID     uuid.UUID  `json:"-"`
	RoomTypeID  uuid.UUID  `json:"room_type_id"`
	RoomNumber  string     `json:"room_number"`
	Status      RoomStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Version     int        `json:"-"`
	RoomType    *RoomType  `json:"room_type,omitempty"` // For embedding on GET requests
}