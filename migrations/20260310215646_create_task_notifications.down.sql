-- Миграция вниз: create_task_notifications

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_delete_task_notifications ON tasks;
DROP FUNCTION IF EXISTS delete_task_notifications_on_deadline_change();
DROP TABLE IF EXISTS task_notifications;
-- +goose StatementEnd