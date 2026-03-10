-- Миграция вверх: create_notification_tables

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notification_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    time_before BIGINT NOT NULL CHECK (time_before > 0),
    enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notification_types_enabled ON notification_types(enabled);

INSERT INTO notification_types (name, time_before) VALUES
    ('3 дня до дедлайна', 259200),
    ('1 день до дедлайна', 86400),
    ('1 час до дедлайна', 3600);
-- +goose StatementEnd