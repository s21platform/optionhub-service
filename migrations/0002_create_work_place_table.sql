-- +goose Up
CREATE TABLE IF NOT EXISTS workplace (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR,
    create_at   TIMESTAMP DEFAULT NOW(),
    is_moderate BOOLEAN   DEFAULT FALSE,
    user_uuid   UUID NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS workplace;
