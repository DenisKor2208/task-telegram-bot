// Package commands
package commands

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btnstart "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands"
)

// StartCommand — структура для команды /start.
type StartCommand struct{}

// Execute — реализация команды /start.
// Отправляет приветственное сообщение с inline-кнопками пользователю.
func (c *StartCommand) Execute(s command.Model, msg messaging.Message) error {
	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := helpers.GetDisplayName(msg)

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTStart, displayName)

	// Отправляем сообщение с inline-кнопками через Telegram-клиент.
	return s.GetTgClient().ShowInlineButtons(text, btnstart.BtnStart, msg.UserID)
}

// NextStep Следующая команда
func (c *StartCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *StartCommand) InputField() string {
	return ""
}
