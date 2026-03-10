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

// GetTaskByID - returns task by tg id
func (ts *TaskStorage) GetTaskByID(ctx context.Context, taskId int64) (*models.Task, error) {
	var task models.Task

	const sqlString = `SELECT * FROM tasks WHERE id = $1`

	// Выполнение запроса на получение данных.
	err := dbutils.Get(ctx, ts.db, &task, sqlString, taskId)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// GetAllTasks - returns all tasks
func (ts *TaskStorage) GetAllTasks(ctx context.Context) ([]*models.Task, error) {
	var tasks []*models.Task

	const sqlString = `SELECT * FROM tasks`

	err := dbutils.Select(ctx, ts.db, &tasks, sqlString)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTasksByStatusID GetAllTasks - returns all tasks
func (ts *TaskStorage) GetTasksByStatusID(ctx context.Context, statusIDs []int) ([]*models.Task, error) {
	if len(statusIDs) == 0 {
		statusIDs = []int{
			StatusInProgress,
			StatusCompleted,
			StatusOverdue,
			StatusClosed,
		}
	}

	var tasks []*models.Task

	const sqlString = `SELECT * FROM tasks WHERE status_id IN (?)`

	// Подготавливаем запрос
	query, args, err := sqlx.In(sqlString, statusIDs)
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

// DeleteTaskByID - deletes task by ID
func (ts *TaskStorage) DeleteTaskByID(ctx context.Context, taskID int64) error {
	const sqlString = `DELETE FROM tasks WHERE id = :id`

	args := map[string]any{"id": taskID}

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
