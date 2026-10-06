-- Миграция вверх: users_tg_id_bigint

-- +goose Up
-- +goose StatementBegin

-- Telegram ID пользователей могут превышать 2^31 (до 52 бит), integer для них мал
ALTER TABLE users ALTER COLUMN tg_id TYPE bigint;
-- +goose StatementEnd
