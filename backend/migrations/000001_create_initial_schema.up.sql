-- Filename: backend/migrations/000001_create_initial_schema.up.sql

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TABLE hotels (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL,
    address     TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hotel_id        UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    email           CITEXT UNIQUE NOT NULL, -- CITEXT is case-insensitive text
    password_hash   BYTEA NOT NULL,
    first_name      TEXT,
    last_name       TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version         INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE roles (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE permissions (
    id   BIGSERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL -- e.g., 'reservations:read', 'users:write'
);

-- Join table for the many-to-many relationship between roles and permissions
CREATE TABLE roles_permissions (
    role_id       BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- Join table for the many-to-many relationship between users and roles
CREATE TABLE users_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Seed with some basic roles and permissions for development
INSERT INTO roles (name) VALUES
('Super Admin'),
('Hotel Manager'),
('Front Desk');

INSERT INTO permissions (code) VALUES
('users:read'),
('users:write'),
('reservations:read'),
('reservations:write'),
('guests:read'),
('guests:write');