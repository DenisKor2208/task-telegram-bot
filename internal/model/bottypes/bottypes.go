// Package bottypes
package bottypes

// Стили (цвета) inline-кнопок Telegram.
const (
	ButtonStylePrimary = "primary" // синяя кнопка
	ButtonStyleSuccess = "success" // зелёная кнопка
	ButtonStyleDanger  = "danger"  // красная кнопка
)

// TgInlineButton - Типы для описания состава кнопок телеграм сообщения.
// Кнопка сообщения.
type TgInlineButton struct {
	DisplayName string
	Value       string
	Style       string
}

// TgRowButtons - Строка с кнопками сообщения.
type TgRowButtons []TgInlineButton
