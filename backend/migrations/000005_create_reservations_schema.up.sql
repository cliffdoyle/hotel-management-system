CREATE TYPE reservation_status AS ENUM (
    'PENDING', 
    'CONFIRMED', 
    'CHECKED_IN', 
    'CHECKED_OUT', 
    'CANCELLED',
    'NO_SHOW'
);

-- This table holds the actual bookings
CREATE TABLE reservations (
    id                    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hotel_id              UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    guest_id              UUID NOT NULL REFERENCES guests(id) ON DELETE CASCADE,
    room_type_id          UUID NOT NULL REFERENCES room_types(id) ON DELETE RESTRICT,
    rate_plan_id          UUID NOT NULL REFERENCES rate_plans(id) ON DELETE RESTRICT,
    
    start_date            DATE NOT NULL,
    end_date              DATE NOT NULL,
    num_adults            INTEGER NOT NULL CHECK (num_adults > 0),
    num_children          INTEGER NOT NULL DEFAULT 0,
    status                reservation_status NOT NULL DEFAULT 'PENDING',
    
    total_cost_cents      INTEGER, -- Storing a snapshot of the price at booking time
    notes                 TEXT,

    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version               INTEGER NOT NULL DEFAULT 1
);

-- This table is our high-performance availability cache
CREATE TABLE inventory_levels (
    hotel_id        UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    room_type_id    UUID NOT NULL REFERENCES room_types(id) ON DELETE CASCADE,
    date            DATE NOT NULL,
    available_rooms INTEGER NOT NULL CHECK (available_rooms >= 0),
    
    PRIMARY KEY (hotel_id, room_type_id, date)
);

-- Add indexes for performance
CREATE INDEX ON reservations(hotel_id);
CREATE INDEX ON reservations(guest_id);
CREATE INDEX ON reservations(room_type_id, start_date, end_date);