package models

import (
	"time"
	"github.com/google/uuid"
)

// RatePlan is a pricing package, e.g., "Breakfast Included", "Non-Refundable".
type RatePlan struct {
	ID                 uuid.UUID `json:"id"`
	HotelID            uuid.UUID `json:"-"`
	Name               string    `json:"name"`
	Description        *string   `json:"description,omitempty"`
	CancellationPolicy *string   `json:"cancellation_policy,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Rate represents the price for a specific room type and rate plan on a single day.
type Rate struct {
	ID           uuid.UUID `json:"id"`
	HotelID      uuid.UUID `json:"-"`
	RoomTypeID   uuid.UUID `json:"room_type_id"`
	RatePlanID   uuid.UUID `json:"rate_plan_id"`
	Date         time.Time `json:"date"`
	PriceCents   int       `json:"price_cents"`
	CreatedAt    time.Time `json:"created_at"`
}