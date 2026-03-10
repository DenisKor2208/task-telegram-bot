package repositories

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TaskNotificationStorage struct {
	db *sqlx.DB
}

func NewTaskNotificationStorage(db *sqlx.DB) *TaskNotificationStorage {
	return &TaskNotificationStorage{db: db}
}

// Insert добавляет запись об отправленном уведомлении. Игнорирует дубликаты.
func (s *TaskNotificationStorage) Insert(ctx context.Context, taskID, notificationTypeID int) error {
	query := `
        INSERT INTO task_notifications (task_id, notification_type_id)
        VALUES ($1, $2)
        ON CONFLICT (task_id, notification_type_id) DO NOTHING
    `
	_, err := s.db.ExecContext(ctx, query, taskID, notificationTypeID)
	if err != nil {
		return fmt.Errorf("Insert task notification: %w", err)
	}
	return nil
}

// GetSentTypesForTask возвращает ID типов уведомлений, уже отправленных для задачи.
func (s *TaskNotificationStorage) GetSentTypesForTask(ctx context.Context, taskID int) ([]int, error) {
	var ids []int
	query := `SELECT notification_type_id FROM task_notifications WHERE task_id = $1`
	err := s.db.SelectContext(ctx, &ids, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("GetSentTypesForTask: %w", err)
	}
	return ids, nil
}
