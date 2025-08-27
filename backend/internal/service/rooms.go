package service

import (
	"context"
	"errors"

	"github.com/cliffdoyle/internal/models"
	"github.com/cliffdoyle/internal/repository"
	validator "github.com/cliffdoyle/internal/validation"
	"github.com/google/uuid"
)

type RoomCreateDTO struct {
	RoomTypeID uuid.UUID         `json:"room_type_id"`
	RoomNumber string            `json:"room_number"`
	Status     models.RoomStatus `json:"status"`
}

type RoomService interface {
	CreateRoom(ctx context.Context, dto RoomCreateDTO, hotelID uuid.UUID) (*models.Room, error)
	GetRoomsForHotel(ctx context.Context, hotelID uuid.UUID) ([]*models.Room, error)
}

type roomService struct {
	roomRepo repository.RoomRepository
}

func NewRoomService(roomRepo repository.RoomRepository) RoomService {
	return &roomService{roomRepo: roomRepo}
}

func (s *roomService) CreateRoom(ctx context.Context, dto RoomCreateDTO, hotelID uuid.UUID) (*models.Room, error) {
	v := validator.New()
	v.Check(dto.RoomNumber != "", "room_number", "must be provided")
	v.Check(dto.RoomTypeID != uuid.Nil, "room_type_id", "must be a valid UUID")

	if !v.Valid() {
		return nil, &validator.ValidationError{Errors: v.Errors}
	}

	room := &models.Room{
		HotelID:    hotelID,
		RoomTypeID: dto.RoomTypeID,
		RoomNumber: dto.RoomNumber,
		Status:     dto.Status,
	}
	// Default status if not provided or invalid
	if room.Status == "" {
		room.Status = models.StatusAvailableClean
	}

	err := s.roomRepo.Insert(ctx, room)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateRoomNumber) {
			v.AddError("room_number", "a room with this number already exists")
			return nil, &validator.ValidationError{Errors: v.Errors}
		}
		return nil, err
	}

	return room, nil
}

func (s *roomService) GetRoomsForHotel(ctx context.Context, hotelID uuid.UUID) ([]*models.Room, error) {
	return s.roomRepo.GetAllForHotel(ctx, hotelID)
}
