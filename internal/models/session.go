package models

// UserSession представляет собой сессию пользователя.
type UserSession struct {
	UserID    int64          `json:"user_id"`    // ID пользователя
	TaskToken string         `json:"task_token"` // Токен задачи
	Data      map[string]any `json:"data"`       // Дополнительные данные сессии
}
