package service

import (
	"context"

	"github.com/cliffdoyle/internal/models"
	"github.com/cliffdoyle/internal/repository"
	validator "github.com/cliffdoyle/internal/validation"
	"github.com/google/uuid"
)

type GuestCreateDTO struct {
	FirstName   string                 `json:"first_name"`
	LastName    string                 `json:"last_name"`
	Email       *string                `json:"email"`
	Phone       *string                `json:"phone"`
	Preferences map[string]interface{} `json:"preferences"`
}

type GuestService interface {
	CreateGuest(ctx context.Context, dto GuestCreateDTO, hotelID uuid.UUID) (*models.Guest, error)
	GetGuestByID(ctx context.Context, hotelID, guestID uuid.UUID) (*models.Guest, error)
	ListGuests(ctx context.Context, hotelID uuid.UUID, filters repository.ListGuestsFilters) ([]*models.Guest, error)
}

type guestService struct {
	guestRepo repository.GuestRepository
}

func NewGuestService(guestRepo repository.GuestRepository) GuestService {
	return &guestService{guestRepo: guestRepo}
}

func (s *guestService) CreateGuest(ctx context.Context, dto GuestCreateDTO, hotelID uuid.UUID) (*models.Guest, error) {
	v := validator.New()
	v.Check(dto.FirstName != "", "first_name", "must be provided")
	v.Check(dto.LastName != "", "last_name", "must be provided")
	if dto.Email != nil {
		v.Check(validator.Matches(*dto.Email, validator.EmailRX), "email", "must be a valid email address")
	}
	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	guest := &models.Guest{
		HotelID:     hotelID,
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Email:       dto.Email,
		Phone:       dto.Phone,
		Preferences: dto.Preferences,
	}

	err := s.guestRepo.Create(ctx, guest)
	if err != nil {
		// No specific business rule errors yet, just pass up repo errors
		return nil, err
	}

	return guest, nil
}

func (s *guestService) GetGuestByID(ctx context.Context, hotelID, guestID uuid.UUID) (*models.Guest, error) {
	return s.guestRepo.GetByID(ctx, hotelID, guestID)
}

func (s *guestService) ListGuests(ctx context.Context, hotelID uuid.UUID, filters repository.ListGuestsFilters) ([]*models.Guest, error) {
	// Default pagination values
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 20
	}

	return s.guestRepo.List(ctx, hotelID, filters)
}
