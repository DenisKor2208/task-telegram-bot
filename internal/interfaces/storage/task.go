package storage

import (
	"context"

	"github.com/DenisKor2208/task-telegram-bot/internal/models"
)

type Task interface {
	GetTaskByID(context.Context, int64) (*models.Task, error)
	CreateTask(context.Context, *models.Task) (bool, error)
	GetAllTasks(context.Context) ([]*models.Task, error)
	GetTasksByStatusID(context.Context, []int) ([]*models.Task, error)
	DeleteTaskByID(context.Context, int64) error
	UpdateTask(context.Context, *models.Task) error
}
