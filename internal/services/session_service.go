package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"fmt"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/callbacktokenpayloadutils"
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

/*********TOKEN**********/

// GenerateToken — создаёт случайный токен из 8 символов
func (s *SessionService) GenerateToken() (string, error) {
	b := make([]byte, 4) // 4 байта = 8 hex-символов
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// CreateCallbackToken — создаёт токен и сохраняет в Redis
func (s *SessionService) CreateCallbackToken(
	ctx context.Context,
	payload callbacktokenpayloadutils.CallbackTokenPayload,
	ttl time.Duration,
) (string, error) {
	token, err := s.GenerateToken()
	if err != nil {
		return "", err
	}

	payload.Expires = time.Now().Add(ttl).Unix()
	key := fmt.Sprintf("cb:%s", token)

	if err := s.storage.Set(ctx, key, payload, ttl); err != nil {
		return "", err
	}

	return token, nil
}

// GetCallbackPayload — получает и удаляет токен (одноразовый)
func (s *SessionService) GetCallbackPayload(
	ctx context.Context,
	token string,
) (*callbacktokenpayloadutils.CallbackTokenPayload, error) {
	key := fmt.Sprintf("cb:%s", token)

	var payload callbacktokenpayloadutils.CallbackTokenPayload
	err := s.storage.GetDel(ctx, key, &payload)
	if err == redis.Nil {
		return nil, fmt.Errorf("token expired or used")
	}
	if err != nil {
		return nil, err
	}

	// Проверяем, не истёк ли токен (если поле Expires задано)
	if time.Now().Unix() > payload.Expires {
		return nil, fmt.Errorf("token expired")
	}

	return &payload, nil
}

// IsCallbackToken — проверяет, похожа ли строка на токен (8 hex-символов)
func (s *SessionService) IsCallbackToken(token string) bool {
	if len(token) != 8 {
		return false
	}
	for _, c := range token {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
