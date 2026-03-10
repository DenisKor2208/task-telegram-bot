// Package service
package service

import (
	"context"

	"github.com/DenisKor2208/task-telegram-bot/internal/models"
)

type Session interface {
	GetSession(context.Context, int64) (*models.UserSession, error)
	GetOrCreateSession(context.Context, int64) (*models.UserSession, error)
	SaveSession(context.Context, *models.UserSession) error
	DeleteSession(context.Context, int64) error
}
