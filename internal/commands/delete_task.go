package commands

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/commands/delete_task"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btndeletetask "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands"
	"github.com/pkg/errors"
)

// DeleteTaskCommand - структура команды /delete_task - "Удалить задачу"
type DeleteTaskCommand struct{}

// Execute — реализация команды /delete_task.
// Отправляет приветственное сообщение с inline-кнопками пользователю.
func (c *DeleteTaskCommand) Execute(s command.Model, msg messaging.Message) error {
	displayName := helpers.GetDisplayName(msg)

	text := fmt.Sprintf(resources.TXTDeleteTask, displayName)

	statuses := []int{
		repositories.StatusInProgress,
		repositories.StatusCompleted,
		repositories.StatusOverdue,
		repositories.StatusClosed,
	}

	tasks, err := s.GetTaskStorage().GetTasksByStatusID(s.GetCtx(), statuses)
	if err != nil {
		return errors.Wrap(err, "Не удалось получить задачи")
	}

	var buttons []bottypes.TgRowButtons
	var additionalButtons []bottypes.TgRowButtons

	for _, task := range tasks {
		deadlineStr := task.Deadline.Format("02.01.2006 15:04")
		if task.Deadline.IsZero() {
			deadlineStr = "Без дедлайна"
		}

		button := bottypes.TgInlineButton{
			DisplayName: fmt.Sprintf("%s (до %s)", task.Description, deadlineStr),
			Value:       fmt.Sprintf("/delete_task_by_id %d", task.ID),
		}
		additionalButtons = append(additionalButtons, bottypes.TgRowButtons{button})
	}

	buttons = append(buttons, additionalButtons...)
	buttons = append(buttons, btndeletetask.BtnDeleteTask...)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)
}

// NextStep Следующая команда
func (c *DeleteTaskCommand) NextStep() command.Command {
	return &delete_task.DeleteTaskByIDCommand{}
}

// InputField Поле для сохранения данных
func (c *DeleteTaskCommand) InputField() string {
	return ""
}
