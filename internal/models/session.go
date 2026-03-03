package models

// UserSession представляет собой сессию пользователя.
// Она хранит состояние FSM и любые временные данные.
type UserSession struct {
	UserID int64                  `json:"user_id"`
	State  string                 `json:"state"` // "ожидает_имя", "ожидает_фамилию" и т.д.
	Data   map[string]interface{} `json:"data"`  // Временные данные, например, "имя"
}
