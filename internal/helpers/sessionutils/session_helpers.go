package sessionutils

import (
	"context"

	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
	"github.com/pkg/errors"
)

// SetSessionTaskToken - устанавливает состояние токена задачи
func SetSessionTaskToken(
	ctx context.Context,
	sessionService types.SessionService,
	userID int64,
	command string,
) (*models.UserSession, error) {

	// 1. Загрузка существующей сессии
	session, _ := sessionService.GetSession(ctx, userID)

	session.State = command

	err := sessionService.SaveSession(ctx, session)
	session, _ = sessionService.GetSession(ctx, userID)
	if err != nil {
		return session, err
	}

	return session, nil
}

// UpdateSessionData загружает текущую сессию пользователя, объединяет (мерджит)
// предоставленные данные в поле Data и сохраняет сессию.
// updateData - это карта, где ключи - это имена полей, а значения - новые данные для обновления.
func UpdateSessionData(
	ctx context.Context,
	sessionService types.SessionService,
	msg types.Message,
	updateData map[string]interface{},
) (*models.UserSession, error) {

	// 1. Загрузка существующей сессии
	session, _ := sessionService.GetSession(ctx, msg.UserID)

	if updateData == nil || len(updateData) == 0 {
		return session, errors.New("данные для обновления не могут быть пустыми")
	}

	// 2. Инициализация Data, если она nil (защита)
	if session.Data == nil {
		session.Data = make(map[string]interface{})
	}

	// 3. Объединение (Merge) новых данных с существующими.
	// Новые ключи перезапишут старые, но остальные останутся.
	for key, value := range updateData {
		session.Data[key] = value
	}

	// 4. Сохранение обновленной сессии обратно в хранилище
	err := sessionService.SaveSession(ctx, session)
	session, _ = sessionService.GetSession(ctx, msg.UserID)
	if err != nil {
		return session, err
	}

	return session, nil
}

// ClearSessionData загружает текущую сессию пользователя, очищает поле Data (устанавливает в пустую карту)
// и сохраняет сессию.
func ClearSessionData(
	ctx context.Context,
	sessionService types.SessionService,
	userID int64,
) (*models.UserSession, error) {

	// 1. Загрузка существующей сессии
	session, _ := sessionService.GetSession(ctx, userID)

	// 2. Очистка поля Data (устанавливаем в пустую карту)
	session.Data = make(map[string]interface{})

	// 3. Сохранение обновленной сессии обратно в хранилище
	err := sessionService.SaveSession(ctx, session)
	session, _ = sessionService.GetSession(ctx, userID)
	if err != nil {
		return session, err
	}

	return session, nil
}
