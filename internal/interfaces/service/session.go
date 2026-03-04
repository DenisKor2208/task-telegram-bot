package service

import (
	"context"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/callbacktokenpayloadutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
)

type Session interface {
	GetSession(context.Context, int64) (*models.UserSession, error)
	GetOrCreateSession(context.Context, int64) (*models.UserSession, error)
	SaveSession(context.Context, *models.UserSession) error
	DeleteSession(context.Context, int64) error

	// Методы для callback-токенов
	CreateCallbackToken(ctx context.Context, payload callbacktokenpayloadutils.CallbackTokenPayload, ttl time.Duration) (string, error)
	GetCallbackPayload(ctx context.Context, token string) (*callbacktokenpayloadutils.CallbackTokenPayload, error)
	IsCallbackToken(token string) bool
}
