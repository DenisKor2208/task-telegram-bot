// Package messaging
package messaging

import "fmt"

// Message Структура сообщения для обработки.
type Message struct {
	Text            string
	Command         string
	Arguments       string
	IsCommand       bool
	UserID          int64
	UserName        string
	UserDisplayName string
	Date            int64
	IsCallback      bool
	CallbackMsgID   string
}

// ValidationError представляет ошибку валидации.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s: %s", e.Field, e.Message)
}

// UserError — ошибка, которую пользователь может исправить сам (например, ошибка ввода).
// Её текст показывается пользователю как есть. Остальные ошибки пользователь видит
// только как общее сообщение, а подробности пишутся в лог.
type UserError struct {
	Text string
}

func (e *UserError) Error() string {
	return e.Text
}

// NewUserError создаёт ошибку с текстом для показа пользователю.
func NewUserError(text string) error {
	return &UserError{Text: text}
}

// Validator определяет контракт для проверки сообщений.
type Validator interface {
	Validate() error
}

// Validate реализует интерфейс Validator для Message.
func (m *Message) Validate() error {
	if m.Command == "" && m.Arguments == "" {
		return ValidationError{Field: "Command/Arguments", Message: "at least one must be provided"}
	}
	return nil
}
