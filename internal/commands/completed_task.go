// Package commands
package commands

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/commands/completed_task"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btncompletedtask "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands"
	"github.com/pkg/errors"
)

// CompletedTaskCommand - структура команды /completed_task - "Выполнить задачу"
type CompletedTaskCommand struct{}

// Execute — реализация команды /completed_task.
func (c *CompletedTaskCommand) Execute(s command.Model, msg messaging.Message) error {

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := helpers.GetDisplayName(msg)

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTCompletedTask, displayName)

	statuses := []int{
		repositories.StatusInProgress,
		//repositories.StatusCompleted,
		repositories.StatusOverdue,
		repositories.StatusClosed,
	}

	tasks, err := s.GetTaskStorage().GetTasksByStatusID(s.GetCtx(), msg.UserID, statuses)
	if err != nil {
		return errors.Wrap(err, "Не удалось получить задачи")
	}

	// Инициализируем кнопки
	var buttons []bottypes.TgRowButtons

	// Инициализируем пустой срез для дополнительных кнопок
	var additionalButtons []bottypes.TgRowButtons

	// Добавляем кнопки для каждой задачи
	for _, task := range tasks {
		deadlineStr := "Без дедлайна"
		if task.Deadline != nil {
			deadlineStr = task.Deadline.Format("02.01.2006 15:04")
		}

		button := bottypes.TgInlineButton{
			DisplayName: fmt.Sprintf("%s (до %s)", task.Description, deadlineStr),
			Value:       fmt.Sprintf("/completed_task_by_id %d", task.ID),
		}

		// Добавляем как новый ряд
		additionalButtons = append(additionalButtons, bottypes.TgRowButtons{button})
	}

	// Формируем общий список кнопок c дополнительными и базовыми кнопками
	buttons = append(buttons, additionalButtons...)
	buttons = append(buttons, btncompletedtask.BtnCompletedTask...)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)

}

// NextStep Следующая команда
func (c *CompletedTaskCommand) NextStep() command.Command {
	return &completed_task.CompletedTaskByIDCommand{}
}

// InputField Поле для сохранения данных
func (c *CompletedTaskCommand) InputField() string {
	return "title"
}
