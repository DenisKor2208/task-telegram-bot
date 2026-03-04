package sessionutils

import (
	"context"
	"maps"

	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/service"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
)

// update — общий helper для обновления сессии с последующим сохранением.
// Принимает функцию, которая модифицирует сессию, и возвращает обновлённую сессию или ошибку.
func update(
	ctx context.Context,
	sessionService service.Session,
	userID int64,
	updateFn func(*models.UserSession),
) (*models.UserSession, error) {
	session, err := sessionService.GetOrCreateSession(ctx, userID)
	if err != nil {
		return nil, err
	}
	updateFn(session)
	err = sessionService.SaveSession(ctx, session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// SetSessionTaskToken устанавливает токен задачи в сессии пользователя.
func SetSessionTaskToken(
	ctx context.Context,
	sessionService service.Session,
	userID int64,
	taskToken string,
) (*models.UserSession, error) {
	return update(ctx, sessionService, userID, func(s *models.UserSession) {
		s.TaskToken = taskToken
	})
}

// UpdateSessionData объединяет (мерджит) переданные данные с полем Data сессии.
// Если updateData пуст, возвращает текущую сессию без изменений (ошибки нет).
func UpdateSessionData(
	ctx context.Context,
	sessionService service.Session,
	userID int64,
	updateData map[string]any,
) (*models.UserSession, error) {
	if len(updateData) == 0 {
		// Нет данных для обновления — просто возвращаем текущую сессию
		return sessionService.GetOrCreateSession(ctx, userID)
	}
	return update(ctx, sessionService, userID, func(s *models.UserSession) {
		if s.Data == nil {
			s.Data = make(map[string]any)
		}
		maps.Copy(s.Data, updateData)
	})
}

// ClearSessionData очищает поле Data сессии (устанавливает пустую карту).
func ClearSessionData(
	ctx context.Context,
	sessionService service.Session,
	userID int64,
) (*models.UserSession, error) {
	return update(ctx, sessionService, userID, func(s *models.UserSession) {
		s.Data = make(map[string]any)
	})
}
