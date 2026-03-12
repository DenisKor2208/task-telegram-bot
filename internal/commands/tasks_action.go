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

// TasksActionCommand — структура для команды /tasks_action.
type TasksActionCommand struct{}

// Execute — реализация команды /tasks_action.
func (c *TasksActionCommand) Execute(s command.Model, msg messaging.Message) error {
	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := helpers.GetDisplayName(msg)

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTStart, displayName)

	// Отправляем сообщение с inline-кнопками через Telegram-клиент.
	return s.GetTgClient().ShowInlineButtons(text, btnstart.BtnTasksAction, msg.UserID)
}

// NextStep Следующая команда
func (c *TasksActionCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *TasksActionCommand) InputField() string {
	return ""
}
