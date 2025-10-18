-- migrations/0002_create_files.sql
-- +goose Up
CREATE TABLE IF NOT EXISTS files (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id      INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT    NOT NULL,
    size         INTEGER NOT NULL,
    bucket       TEXT    NOT NULL,
    object_name  TEXT    NOT NULL,
    uploaded_at  DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_files_user_id ON files(user_id);

-- +goose Down
DROP TABLE IF EXISTS files;
