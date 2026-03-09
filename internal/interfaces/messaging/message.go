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
