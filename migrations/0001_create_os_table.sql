-- +goose Up
CREATE TABLE IF NOT EXISTS os (
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    create_at timestamp,
    is_moderate boolean default false,
    user_uuid varchar
);



-- +goose Down
DROP TABLE IF EXISTS os;