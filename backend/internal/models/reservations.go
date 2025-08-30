package models

import (
	"time"

	"github.com/google/uuid"
)

type ReservationStatus string

const (
	StatusPending    ReservationStatus = "PENDING"
	StatusConfirmed  ReservationStatus = "CONFIRMED"
	StatusCheckedIn  ReservationStatus = "CHECKED_IN"
	StatusCheckedOut ReservationStatus = "CHECKED_OUT"
	StatusCancelled  ReservationStatus = "CANCELLED"
	StatusNoShow     ReservationStatus = "NO_SHOW"
)

type Reservation struct {
	ID             uuid.UUID         `json:"id"`
	HotelID        uuid.UUID         `json:"-"`
	GuestID        uuid.UUID         `json:"guest_id"`
	RoomTypeID     uuid.UUID         `json:"room_type_id"`
	RatePlanID     uuid.UUID         `json:"rate_plan_id"`
	StartDate      time.Time         `json:"start_date"`
	EndDate        time.Time         `json:"end_date"`
	NumAdults      int               `json:"num_adults"`
	NumChildren    int               `json:"num_children"`
	Status         ReservationStatus `json:"status"`
	TotalCostCents int               `json:"total_cost_cents"`
	Notes          *string           `json:"notes,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	Version        int               `json:"-"`
	Guest          *Guest            `json:"guest,omitempty"`
}
