CREATE TABLE guests (
    id                    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hotel_id              UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    first_name            TEXT NOT NULL,
    last_name             TEXT NOT NULL,
    email                 CITEXT,
    phone                 TEXT,
    
    -- Loyalty program fields
    loyalty_member_number TEXT,
    loyalty_tier          TEXT,
    loyalty_points        INTEGER DEFAULT 0,

    -- Flexible columns for future needs
    preferences           JSONB, -- e.g., {"room_preference": "high_floor", "pillow_type": "foam"}
    metadata              JSONB, -- For internal or integration-specific data

    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version               INTEGER NOT NULL DEFAULT 1,

    -- A guest's email must be unique within a hotel, but can be NULL
    UNIQUE (hotel_id, email)
);

CREATE TABLE guest_communication_logs (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    guest_id          UUID NOT NULL REFERENCES guests(id) ON DELETE CASCADE,
    sent_by_user_id   UUID REFERENCES users(id) ON DELETE SET NULL, -- Can be sent by system (NULL) or user
    channel           TEXT NOT NULL, -- e.g., 'EMAIL', 'SMS', 'INTERNAL_NOTE'
    message           TEXT NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX ON guests(hotel_id);
CREATE INDEX ON guests(hotel_id, lower(email));
-- A GIN index allows for fast searches inside the JSONB columns if needed later
CREATE INDEX ON guests USING GIN (preferences); 
-- A trigram index is great for "fuzzy" name searching
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX ON guests USING GIN (first_name gin_trgm_ops, last_name gin_trgm_ops);

CREATE INDEX ON guest_communication_logs(guest_id);