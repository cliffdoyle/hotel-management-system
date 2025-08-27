CREATE TABLE room_types (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hotel_id    UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    description TEXT,
    capacity    INTEGER NOT NULL CHECK (capacity > 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- A room type's name must be unique within its hotel
    UNIQUE (hotel_id, name)
);

CREATE TYPE room_status AS ENUM (
    'AVAILABLE_CLEAN',
    'AVAILABLE_DIRTY',
    'OCCUPIED',
    'OUT_OF_SERVICE'
);

CREATE TABLE rooms (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hotel_id     UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    room_type_id UUID NOT NULL REFERENCES room_types(id) ON DELETE RESTRICT,
    room_number  TEXT NOT NULL,
    status       room_status NOT NULL DEFAULT 'AVAILABLE_DIRTY',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version      INTEGER NOT NULL DEFAULT 1,

    -- A room number must be unique within its hotel
    UNIQUE (hotel_id, room_number)
);

-- Add indexes for performance on foreign keys
CREATE INDEX ON room_types(hotel_id);
CREATE INDEX ON rooms(hotel_id);
CREATE INDEX ON rooms(room_type_id);