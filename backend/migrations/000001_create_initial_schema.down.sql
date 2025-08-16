-- Filename: backend/migrations/000001_create_initial_schema.down.sql

DROP TABLE IF EXISTS users_roles;
DROP TABLE IF EXISTS roles_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS hotels;

DROP EXTENSION IF EXISTS "uuid-ossp";