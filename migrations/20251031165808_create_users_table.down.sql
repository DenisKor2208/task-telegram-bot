-- Миграция вниз: create_users_table

-- +goose Down
-- +goose StatementBegin

-- Drop trigger and function first
drop trigger if exists update_users_updated_at on users;
drop function if exists update_updated_at_column();

-- Drop index
drop index if exists users_tg_id_idx;

-- Drop table
drop table if exists users;
-- +goose StatementEnd