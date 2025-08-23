-- +goose Up
CREATE TABLE IF NOT EXISTS attributes
(
    attribute_id      SERIAL PRIMARY KEY,
    name              TEXT        NOT NULL, -- машинное имя ("city", "birth_date")
    type              VARCHAR(50) NOT NULL, -- "string", "number", "date", "option", "multiselect"
    description       TEXT,
    allowed_operators JSONB,
    created_at        TIMESTAMP DEFAULT now(),
    updated_at        TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS options
(
    option_id    SERIAL PRIMARY KEY,
    attribute_id INT  NOT NULL,
    label        TEXT NOT NULL,                      -- "Челябинск", "Москва"
    parent_id    INT,                                -- для иерархии (страна → регион → город)
    sort_order   INT       DEFAULT 0,
    created_at   TIMESTAMP DEFAULT now(),
    updated_at   TIMESTAMP DEFAULT now(),

    CONSTRAINT fk_options_attribute
        FOREIGN KEY (attribute_id)
            REFERENCES attributes (attribute_id)
            ON DELETE CASCADE,
    
    CONSTRAINT fk_options_parent
        FOREIGN KEY (parent_id)
            REFERENCES options (option_id)
);

CREATE TABLE IF NOT EXISTS entity_attributes
(
    entity_attribute_id SERIAL PRIMARY KEY,
    entity_type         VARCHAR(50) NOT NULL, -- "user", "event", "article" и т.д.
    attribute_id        INT         NOT NULL,
    label               TEXT        NOT NULL, -- UI-лейбл для фронта ("город проживания", "город мероприятия")
    is_required         BOOLEAN   DEFAULT FALSE,
    order_index         INT       DEFAULT 0,  -- порядок отображения в форме
    visibility_rules    JSONB,                -- гибкие правила видимости
    created_at          TIMESTAMP DEFAULT now(),
    updated_at          TIMESTAMP DEFAULT now(),

    CONSTRAINT fk_entity_attributes_attribute
        FOREIGN KEY (attribute_id)
            REFERENCES attributes (attribute_id)
            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS entity_attributes;
DROP TABLE IF EXISTS options;
DROP TABLE IF EXISTS attributes;
