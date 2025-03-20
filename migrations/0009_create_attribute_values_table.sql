-- +goose Up
CREATE TABLE attribute_values
(
    id           SERIAL PRIMARY KEY,
    attribute_id INT REFERENCES attributes (id) ON DELETE CASCADE,
    value        TEXT NOT NULL,
    parent_id    INT REFERENCES attribute_values (id),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS attribute_values;
