// Package services
package services

import (
	"context"

	"fmt"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// SessionService который работает с сессиями
type SessionService struct {
	storage    *repositories.RedisStorage
	sessionTTL time.Duration // Дефолтное TTL для сессий
}

// NewSessionService — конструктор
func NewSessionService(storage *repositories.RedisStorage, defaultTTL time.Duration) *SessionService {
	return &SessionService{
		storage:    storage,
		sessionTTL: defaultTTL,
	}
}

// createSessionKey — хелпер для генерации ключа
func createSessionKey(userID int64) string {
	return fmt.Sprintf("session:%d", userID)
}

// SaveSession — сохраняет сессию в Redis
func (s *SessionService) SaveSession(ctx context.Context, session *models.UserSession) error {
	key := createSessionKey(session.UserID)
	err := s.storage.Set(ctx, key, session, s.sessionTTL)
	if err != nil {
		return fmt.Errorf("save session for user %d: %w", session.UserID, err)
	}
	return nil
}

// GetSession — получает сессию из Redis
func (s *SessionService) GetSession(ctx context.Context, userID int64) (*models.UserSession, error) {
	key := createSessionKey(userID)

	var session models.UserSession

	err := s.storage.Get(ctx, key, &session)
	if err != nil {
		return nil, fmt.Errorf("get session for user %d: %w", userID, err)
	}
	return &session, nil

}

// GetOrCreateSession — получает сессию или создаёт новую, если её нет.
// Новая сессия сразу сохраняется в Redis с дефолтным TTL.
func (s *SessionService) GetOrCreateSession(ctx context.Context, userID int64) (*models.UserSession, error) {
	session, err := s.GetSession(ctx, userID)
	if err == nil {
		return session, nil
	}
	if !errors.Is(err, redis.Nil) {
		return nil, err // реальная ошибка
	}
	// Сессия не найдена — создаём новую
	newSession := &models.UserSession{
		UserID:    userID,
		TaskToken: "",
		Data:      make(map[string]any),
	}
	if err := s.SaveSession(ctx, newSession); err != nil {
		return nil, fmt.Errorf("failed to save new session for user %d: %w", userID, err)
	}
	return newSession, nil
}

// DeleteSession — удаляет сессию (например, при /stop)
func (s *SessionService) DeleteSession(ctx context.Context, userID int64) error {
	key := createSessionKey(userID)
	_, err := s.storage.Del(ctx, key) // `Del` возвращает (int64, error)
	return err
}
