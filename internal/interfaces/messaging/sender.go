package messaging

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// Sender Интерфейс для работы с сообщениями.
type Sender interface {
	SendMessage(text string, userID int64) error
	ShowInlineButtons(text string, buttons []bottypes.TgRowButtons, userID int64) error
	SendMessageWithTemplate(text string, textTemplate string, userID int64) error
}
