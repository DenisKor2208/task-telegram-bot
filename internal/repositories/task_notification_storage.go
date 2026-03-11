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

// GetSentMapForTasks возвращает для каждого task_id список notification_type_id,
// которые уже были отправлены. Возвращает мапу, где ключ — task_id, значение — срез типов.
func (s *TaskNotificationStorage) GetSentMapForTasks(ctx context.Context, taskIDs []int) (map[int][]int, error) {
	if len(taskIDs) == 0 {
		return map[int][]int{}, nil
	}

	// Используем ANY($1) для передачи слайса
	query := `SELECT task_id, notification_type_id FROM task_notifications WHERE task_id = ANY($1)`
	rows, err := s.db.QueryContext(ctx, query, taskIDs)
	if err != nil {
		return nil, fmt.Errorf("GetSentMapForTasks query: %w", err)
	}
	defer rows.Close()

	result := make(map[int][]int)
	for rows.Next() {
		var taskID, notifTypeID int
		if err := rows.Scan(&taskID, &notifTypeID); err != nil {
			return nil, fmt.Errorf("GetSentMapForTasks scan: %w", err)
		}
		result[taskID] = append(result[taskID], notifTypeID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetSentMapForTasks rows: %w", err)
	}
	return result, nil
}
