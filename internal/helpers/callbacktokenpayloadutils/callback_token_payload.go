// Package callbacktokenpayloadutils предоставляет утилиты для работы с полезной нагрузкой callback-токена.
package callbacktokenpayloadutils

// CallbackTokenPayload — данные для callback-токена
type CallbackTokenPayload struct {
	Action  string `json:"a"`           // "edit", "delete", "complete"
	Entity  string `json:"e"`           // "task" -
	TaskID  int64  `json:"t"`           // ID задачи
	UserID  int64  `json:"u"`           // Кто нажал (для проверки прав)
	Field   string `json:"f,omitempty"` // Какое поле редактируем
	Expires int64  `json:"exp"`         // Когда истечёт (unix timestamp)
}
