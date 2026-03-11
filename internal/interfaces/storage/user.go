package storage

import (
	"context"

	"github.com/DenisKor2208/task-telegram-bot/internal/models"
)

type User interface {
	GetUserByTgID(context.Context, int) (*models.User, error)
	GetAllUsers(context.Context, int, int) ([]*models.User, error)
	GetUserByID(context.Context, int) (*models.User, error)
	CreateUser(context.Context, *models.User) (*models.User, error)
	UpdateUserTimezone(context.Context, int, string) error
}
