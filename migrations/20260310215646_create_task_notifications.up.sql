-- Миграция вверх: create_task_notifications

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS task_notifications (
    id SERIAL PRIMARY KEY,
    task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    notification_type_id INT NOT NULL REFERENCES notification_types(id) ON DELETE CASCADE,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(task_id, notification_type_id)
);

CREATE INDEX IF NOT EXISTS idx_task_notifications_task_id ON task_notifications(task_id);

-- Функция, удаляющая уведомления задачи при изменении deadline
CREATE OR REPLACE FUNCTION delete_task_notifications_on_deadline_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.deadline IS DISTINCT FROM NEW.deadline THEN
        DELETE FROM task_notifications WHERE task_id = NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер
DROP TRIGGER IF EXISTS trigger_delete_task_notifications ON tasks;
CREATE TRIGGER trigger_delete_task_notifications
    BEFORE UPDATE OF deadline ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION delete_task_notifications_on_deadline_change();
-- +goose StatementEnd