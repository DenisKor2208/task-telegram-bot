-- Миграция вниз: tasks
-- +goose Down

-- Drop trigger
drop trigger if exists update_tasks_updated_at on tasks;

-- Drop indexes
drop index if exists tasks_user_id_idx;
drop index if exists tasks_status_id_idx;
drop index if exists tasks_deadline_idx;
drop index if exists tasks_priority_idx;

-- Drop table
drop table if exists tasks;
-- +goose StatementEnd
