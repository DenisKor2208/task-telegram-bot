// Package bottypes
package bottypes

// TgInlineButton - Типы для описания состава кнопок телеграм сообщения.
// Кнопка сообщения.
type TgInlineButton struct {
	DisplayName string
	Value       string
}

// TgRowButtons - Строка с кнопками сообщения.
type TgRowButtons []TgInlineButton
