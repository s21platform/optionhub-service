-- +goose Up
CREATE TABLE attributes
(
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    type       TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS attributes;
