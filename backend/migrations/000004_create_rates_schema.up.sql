CREATE TABLE rate_plans (
    id                    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hotel_id              UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    name                  TEXT NOT NULL,
    description           TEXT,
    cancellation_policy   TEXT, -- For storing cancellation rules text
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- A rate plan's name must be unique within its hotel
    UNIQUE(hotel_id, name)
);

CREATE TABLE rates (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hotel_id        UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    room_type_id    UUID NOT NULL REFERENCES room_types(id) ON DELETE CASCADE,
    rate_plan_id    UUID NOT NULL REFERENCES rate_plans(id) ON DELETE CASCADE,
    
    date            DATE NOT NULL,
    price_cents     INTEGER NOT NULL CHECK (price_cents >= 0),
    
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- We can only have one price for a given room type, rate plan, on a specific day
    UNIQUE(hotel_id, room_type_id, rate_plan_id, date)
);

-- Indexes for fast lookups, essential for the pricing engine
CREATE INDEX ON rate_plans(hotel_id);
CREATE INDEX ON rates(hotel_id, room_type_id, rate_plan_id, date);