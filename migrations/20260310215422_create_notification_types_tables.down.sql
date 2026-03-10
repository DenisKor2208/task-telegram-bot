-- Миграция вниз: create_notification_tables

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_types;
-- +goose StatementEnd