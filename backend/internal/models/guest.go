package models

import (
	"time"
	"github.com/google/uuid"
)

type Guest struct {
	ID                  uuid.UUID              `json:"id"`
	HotelID             uuid.UUID              `json:"-"`
	FirstName           string                 `json:"first_name"`
	LastName            string                 `json:"last_name"`
	Email               *string                `json:"email,omitempty"`
	Phone               *string                `json:"phone,omitempty"`
	LoyaltyMemberNumber *string                `json:"loyalty_member_number,omitempty"`
	LoyaltyTier         *string                `json:"loyalty_tier,omitempty"`
	LoyaltyPoints       int                    `json:"loyalty_points"`
	Preferences         map[string]interface{} `json:"preferences,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	Version             int                    `json:"-"`
}

type GuestCommunicationLog struct {
	ID            uuid.UUID  `json:"id"`
	GuestID       uuid.UUID  `json:"guest_id"`
	SentByUserID  *uuid.UUID `json:"sent_by_user_id,omitempty"`
	Channel       string     `json:"channel"`
	Message       string     `json:"message"`
	CreatedAt     time.Time  `json:"created_at"`
}