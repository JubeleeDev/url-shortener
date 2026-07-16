-- +goose Up
CREATE TABLE links (
    code TEXT PRIMARY KEY,
    original_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now());


-- +goose Down
DROP TABLE IF EXISTS links;
