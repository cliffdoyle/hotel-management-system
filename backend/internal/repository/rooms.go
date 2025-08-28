package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/cliffdoyle/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDuplicateRoomNumber = errors.New("a room with this number already exists in this hotel")
)

type RoomRepository interface {
	Insert(ctx context.Context, room *models.Room) error
	GetAllForHotel(ctx context.Context, hotelID uuid.UUID) ([]*models.Room, error)
	// GetByID, Update, Delete can be added here later
}

type roomRepository struct {
	db *pgxpool.Pool
}

func NewRoomRepository(db *pgxpool.Pool) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Insert(ctx context.Context, room *models.Room) error {
	query := `
		INSERT INTO rooms (hotel_id, room_type_id, room_number, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at, version`

	args := []any{room.HotelID, room.RoomTypeID, room.RoomNumber, room.Status}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, query, args...).Scan(&room.ID, &room.CreatedAt, &room.UpdatedAt, &room.Version)
	if err != nil {
		if strings.Contains(err.Error(), "rooms_hotel_id_room_number_key") {
			return ErrDuplicateRoomNumber
		}
		return err
	}
	return nil
}

func (r *roomRepository) GetAllForHotel(ctx context.Context, hotelID uuid.UUID) ([]*models.Room, error) {
	query := `
		SELECT 
			r.id, r.hotel_id, r.room_type_id, r.room_number, r.status, 
			r.created_at, r.updated_at, r.version,
			rt.id, rt.name, rt.description, rt.capacity
		FROM rooms r
		INNER JOIN room_types rt ON r.room_type_id = rt.id
		WHERE r.hotel_id = $1
		ORDER BY r.room_number ASC`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, query, hotelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		var room models.Room
		var roomType models.RoomType

		err := rows.Scan(
			&room.ID, &room.HotelID, &room.RoomTypeID, &room.RoomNumber, &room.Status,
			&room.CreatedAt, &room.UpdatedAt, &room.Version,
			&roomType.ID, &roomType.Name, &roomType.Description, &roomType.Capacity,
		)
		if err != nil {
			return nil, err
		}

		room.RoomType = &roomType
		rooms = append(rooms, &room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}
