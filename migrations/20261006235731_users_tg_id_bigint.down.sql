-- Миграция вниз: users_tg_id_bigint

-- +goose Down
-- +goose StatementBegin

-- Откат невозможен, если в таблице уже есть tg_id больше 2147483647 (integer out of range)
ALTER TABLE users ALTER COLUMN tg_id TYPE integer;
-- +goose StatementEnd
