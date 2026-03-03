package commands

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
	"github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands"
)

const (
	COMMAND_VIEW_TASKS = "/view_tasks"
)

// ViewTasksCommand — структура для команды /view_tasks - "Просмотреть задачи".
type ViewTasksCommand struct{}

// Execute — реализация команды /view_tasks.
// Отправляет приветственное сообщение с inline-кнопками пользователю.
func (c *ViewTasksCommand) Execute(s types.Model, msg types.Message, session *models.UserSession) error {
	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := msg.UserDisplayName
	if len(displayName) == 0 {
		displayName = msg.UserName
	}

	text := fmt.Sprintf(resources.TXTViewTasks, displayName)

	// Отправляем сообщение с inline-кнопками через Telegram-клиент.
	return s.GetTgClient().ShowInlineButtons(text, commands.BtnViewTasks, msg.UserID)
}
