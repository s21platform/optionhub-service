-- +goose Up
CREATE TABLE IF NOT EXISTS attribute_request
(
    id           SERIAL PRIMARY KEY,
    attribute_id INT REFERENCES attributes (id) ON DELETE CASCADE,
    value        TEXT NOT NULL,
    user_uuid    UUID NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS attribute_request;
