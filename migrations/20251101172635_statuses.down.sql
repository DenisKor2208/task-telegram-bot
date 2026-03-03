-- Миграция вниз: statuses

-- +goose Down
-- +goose StatementBegin

-- Drop trigger
drop trigger if exists update_statuses_updated_at on statuses;

-- Drop index
drop index if exists statuses_name_idx;

-- Drop table
drop table if exists statuses;
-- +goose StatementEnd