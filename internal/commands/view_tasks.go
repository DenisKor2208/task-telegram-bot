package commands

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btnviewtasks "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands"
)

// ViewTasksCommand — структура для команды /view_tasks - "Просмотреть задачи".
type ViewTasksCommand struct{}

// Execute — реализация команды /view_tasks.
// Отправляет приветственное сообщение с inline-кнопками пользователю.
func (c *ViewTasksCommand) Execute(s command.Model, msg messaging.Message) error {
	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := helpers.GetDisplayName(msg)

	text := fmt.Sprintf(resources.TXTViewTasks, displayName)

	// Отправляем сообщение с inline-кнопками через Telegram-клиент.
	return s.GetTgClient().ShowInlineButtons(text, btnviewtasks.BtnViewTasks, msg.UserID)
}

// NextStep Следующая команда
func (c *ViewTasksCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *ViewTasksCommand) InputField() string {
	return ""
}
