package storage

import (
	"context"

	"github.com/DenisKor2208/task-telegram-bot/internal/models"
)

type User interface {
	GetUserByTgID(context.Context, int) (*models.User, error)
	CreateUser(context.Context, *models.User) (*models.User, error)
}
