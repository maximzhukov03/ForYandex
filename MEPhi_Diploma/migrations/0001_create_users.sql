-- migrations/0001_users.sql

-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id             INTEGER    PRIMARY KEY AUTOINCREMENT,
    email          TEXT       NOT NULL UNIQUE,
    password_hash  TEXT       NOT NULL,
    created_at     DATETIME   NOT NULL,
    role           TEXT       NOT NULL DEFAULT 'user'
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- +goose Down
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;