package commands

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
	"github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands"
)

const (
	COMMAND_START = "/start"
)

// StartCommand — структура для команды /start.
type StartCommand struct{}

// Execute — реализация команды /start.
// Отправляет приветственное сообщение с inline-кнопками пользователю.
func (c *StartCommand) Execute(s types.Model, msg types.Message, session *models.UserSession) error {
	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := msg.UserDisplayName
	if len(displayName) == 0 {
		displayName = msg.UserName
	}

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTStart, displayName)

	// Отправляем сообщение с inline-кнопками через Telegram-клиент.
	return s.GetTgClient().ShowInlineButtons(text, commands.BtnStart, msg.UserID)
}
