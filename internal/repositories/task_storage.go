package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/dbutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/jmoiron/sqlx"
)

// TaskStorage - task storage
type TaskStorage struct {
	db *sqlx.DB
}

// NewTaskStorage returns new instance of TaskStorage
func NewTaskStorage(db *sqlx.DB) *TaskStorage {
	return &TaskStorage{db: db}
}

// GetTaskByID возвращает задачу по ID, только если она принадлежит пользователю с Telegram ID tgID.
// Для чужой или несуществующей задачи возвращает ошибку (sql.ErrNoRows).
func (ts *TaskStorage) GetTaskByID(ctx context.Context, tgID int64, taskID int64) (*models.Task, error) {
	var task models.Task

	const sqlString = `
		SELECT t.* FROM tasks t
		JOIN users u ON u.id = t.user_id
		WHERE t.id = $1 AND u.tg_id = $2`

	// Выполнение запроса на получение данных.
	err := dbutils.Get(ctx, ts.db, &task, sqlString, taskID, tgID)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// GetAllTasks возвращает все задачи пользователя с Telegram ID tgID.
func (ts *TaskStorage) GetAllTasks(ctx context.Context, tgID int64) ([]*models.Task, error) {
	var tasks []*models.Task

	const sqlString = `
		SELECT t.* FROM tasks t
		JOIN users u ON u.id = t.user_id
		WHERE u.tg_id = $1
		ORDER BY t.deadline NULLS LAST, t.created_at`

	err := dbutils.Select(ctx, ts.db, &tasks, sqlString, tgID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTasksByStatusID возвращает задачи пользователя с Telegram ID tgID в указанных статусах.
// Если статусы не переданы, возвращаются задачи во всех статусах.
func (ts *TaskStorage) GetTasksByStatusID(ctx context.Context, tgID int64, statusIDs []int) ([]*models.Task, error) {
	if len(statusIDs) == 0 {
		statusIDs = []int{
			StatusInProgress,
			StatusCompleted,
			StatusOverdue,
			StatusClosed,
		}
	}

	var tasks []*models.Task

	const sqlString = `
		SELECT t.* FROM tasks t
		JOIN users u ON u.id = t.user_id
		WHERE u.tg_id = ? AND t.status_id IN (?)
		ORDER BY t.deadline NULLS LAST, t.created_at`

	// Подготавливаем запрос
	query, args, err := sqlx.In(sqlString, tgID, statusIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	// Преобразуем ? -> $1, $2, ... для PostgreSQL
	query = ts.db.Rebind(query)

	err = dbutils.Select(ctx, ts.db, &tasks, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить задачи по IDs: %w", err)
	}

	return tasks, nil
}

// CreateTask - creates task
func (ts *TaskStorage) CreateTask(ctx context.Context, task *models.Task) (bool, error) {

	const sqlString = `
		INSERT INTO tasks (description, deadline, priority, status_id, user_id, created_at, updated_at)
		VALUES (:description, :deadline, :priority, :status_id, :user_id, :created_at, :updated_at)
	`

	now := time.Now()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = now
	}

	if task.UpdatedAt.IsZero() {
		task.UpdatedAt = now
	}

	_, err := dbutils.NamedExec(ctx, ts.db, sqlString, task)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteTaskByID удаляет задачу по ID, только если она принадлежит пользователю с Telegram ID tgID.
// Для чужой или несуществующей задачи возвращает ошибку «не найдена».
func (ts *TaskStorage) DeleteTaskByID(ctx context.Context, tgID int64, taskID int64) error {
	const sqlString = `
		DELETE FROM tasks t
		USING users u
		WHERE t.id = :id AND t.user_id = u.id AND u.tg_id = :tg_id`

	args := map[string]any{"id": taskID, "tg_id": tgID}

	// Выполнение запроса на получение данных.
	result, err := dbutils.NamedExec(ctx, ts.db, sqlString, args)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача с ID %d не найдена", taskID)
	}

	return nil
}

// UpdateTask - updates task by ID
func (ts *TaskStorage) UpdateTask(ctx context.Context, task *models.Task) error {
	if task.ID <= 0 {
		return fmt.Errorf("invalid task ID: %d", task.ID)
	}

	// Устанавливаем время обновления
	task.UpdatedAt = time.Now()

	const sqlString = `
		UPDATE tasks 
		SET description = :description, 
		    deadline = :deadline, 
		    priority = :priority, 
		    status_id = :status_id, 
		    user_id = :user_id, 
		    updated_at = :updated_at 
		WHERE id = :id
	`

	result, err := dbutils.NamedExec(ctx, ts.db, sqlString, task)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача с ID %d не найдена", task.ID)
	}

	return nil
}

// UpdateOverdueTasks переводит в "Просрочено" задачи в статусе "В процессе", у которых дедлайн в прошлом.
// Задачи в статусах "Выполнено", "Просрочено" и "Завершено" не затрагиваются.
// Возвращает количество обновлённых строк.
func (ts *TaskStorage) UpdateOverdueTasks(ctx context.Context) (int64, error) {
	query := `UPDATE tasks
						SET status_id = $1, updated_at = NOW()
            WHERE deadline < NOW() AND status_id = $2`
	result, err := ts.db.ExecContext(ctx, query, StatusOverdue, StatusInProgress)
	if err != nil {
		return 0, fmt.Errorf("failed to update overdue tasks: %w", err)
	}
	return result.RowsAffected()
}

// GetTasksForNotification возвращает только задачи, которые потенциально могут требовать уведомлений (активные, с будущим дедлайном)
func (ts *TaskStorage) GetTasksForNotification(ctx context.Context) ([]*models.Task, error) {
	var tasks []*models.Task
	query := `SELECT * FROM tasks 
              WHERE deadline > NOW() 
                AND status_id NOT IN ($1, $2, $3) 
              ORDER BY deadline`
	err := dbutils.Select(ctx, ts.db, &tasks, query, StatusOverdue, StatusClosed, StatusCompleted)
	if err != nil {
		return nil, fmt.Errorf("GetTasksForNotification: %w", err)
	}
	return tasks, nil
}
