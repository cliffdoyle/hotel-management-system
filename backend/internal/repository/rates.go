package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cliffdoyle/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RateRepository interface {
	// For admin: Creates multiple daily rate entries from a date range
	BulkInsert(ctx context.Context, rates []*models.Rate) (int64, error)
	// For pricing engine: Finds all daily rates for a specific stay query
	FindForStay(ctx context.Context, hotelID, roomTypeID, ratePlanID uuid.UUID, startDate, endDate time.Time) ([]*models.Rate, error)
}

type rateRepository struct {
	db *pgxpool.Pool
}

func NewRateRepository(db *pgxpool.Pool) RateRepository {
	return &rateRepository{db: db}
}

func (r *rateRepository) BulkInsert(ctx context.Context, rates []*models.Rate) (int64, error) {
	// pgx CopyFrom is highly efficient for bulk inserts.
	// We define the columns we're inserting into.
	columns := []string{"hotel_id", "room_type_id", "rate_plan_id", "date", "price_cents"}

	// Create a slice of slices, where each inner slice represents a row.
	rows := make([][]interface{}, len(rates))
	for i, rate := range rates {
		rows[i] = []interface{}{rate.HotelID, rate.RoomTypeID, rate.RatePlanID, rate.Date, rate.PriceCents}
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Use the CopyFrom helper
	copyCount, err := r.db.CopyFrom(
		ctx,
		pgx.Identifier{"rates"},
		columns,
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		if strings.Contains(err.Error(), "rates_hotel_id_room_type_id_rate_plan_id_date_key") {
			return 0, fmt.Errorf("one or more rates for the given dates already exist")
		}
		return 0, err
	}

	return copyCount, nil
}

func (r *rateRepository) FindForStay(ctx context.Context, hotelID, roomTypeID, ratePlanID uuid.UUID, startDate, endDate time.Time) ([]*models.Rate, error) {
	// Note: endDate is exclusive in this query. A 1-night stay from 2024-01-01 to 2024-01-02
	// means we only need the rate for the night of the 1st.
	query := `
		SELECT id, date, price_cents
		FROM rates
		WHERE hotel_id = $1
		  AND room_type_id = $2
		  AND rate_plan_id = $3
		  AND date >= $4 AND date < $5
		ORDER BY date`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, query, hotelID, roomTypeID, ratePlanID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []*models.Rate
	for rows.Next() {
		var rate models.Rate
		err := rows.Scan(&rate.ID, &rate.Date, &rate.PriceCents)
		if err != nil {
			return nil, err
		}
		rates = append(rates, &rate)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rates, nil
}

// NOTE: A RatePlanRepository would also be created in a similar fashion.
