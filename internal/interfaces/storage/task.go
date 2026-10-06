package storage

import (
	"context"

	"github.com/DenisKor2208/task-telegram-bot/internal/models"
)

// Task — хранилище задач. Методы с параметром tgID работают только
// с задачами пользователя с этим Telegram ID.
type Task interface {
	GetTaskByID(ctx context.Context, tgID int64, taskID int64) (*models.Task, error)
	CreateTask(context.Context, *models.Task) (bool, error)
	GetAllTasks(ctx context.Context, tgID int64) ([]*models.Task, error)
	GetTasksByStatusID(ctx context.Context, tgID int64, statusIDs []int) ([]*models.Task, error)
	DeleteTaskByID(ctx context.Context, tgID int64, taskID int64) error
	UpdateTask(context.Context, *models.Task) error
}
