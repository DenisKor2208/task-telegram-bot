package commands

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/commands/closed_task"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btnclosedtask "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands"
	"github.com/pkg/errors"
)

// ClosedTaskCommand - структура команды /closed_task - "Завершить задачу"
type ClosedTaskCommand struct{}

// Execute — реализация команды /closed_task.
// Отправляет приветственное сообщение с inline-кнопками пользователю.
func (c *ClosedTaskCommand) Execute(s command.Model, msg messaging.Message) error {

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := helpers.GetDisplayName(msg)

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTClosedTask, displayName)

	statuses := []int{
		repositories.STATUS_IN_PROGRESS,
		repositories.STATUS_COMPLETED,
		repositories.STATUS_OVERDUE,
		// repositories.STATUS_CLOSED,
	}

	tasks, err := s.GetTaskStorage().GetTasksByStatusID(s.GetCtx(), statuses)
	if err != nil {
		return errors.Wrap(err, "Не удалось получить задачи")
	}

	// Инициализируем кнопки
	var buttons []bottypes.TgRowButtons

	// Инициализируем пустой срез для дополнительных кнопок
	var additionalButtons []bottypes.TgRowButtons

	// Добавляем кнопки для каждой задачи
	for _, task := range tasks {
		deadlineStr := task.Deadline.Format("02.01.2006 15:04")
		if task.Deadline.IsZero() {
			deadlineStr = "Без дедлайна"
		}

		button := bottypes.TgInlineButton{
			DisplayName: fmt.Sprintf("%s (до %s)", task.Description, deadlineStr),
			Value:       fmt.Sprintf("/closed_task_by_id %d", task.ID),
		}

		// Добавляем как новый ряд
		additionalButtons = append(additionalButtons, bottypes.TgRowButtons{button})
	}

	// Формируем общий список кнопок с дополнительными и базовыми кнопками
	buttons = append(buttons, additionalButtons...)
	buttons = append(buttons, btnclosedtask.BtnClosedTask...)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)

}

// NextStep Следующая команда
func (c *ClosedTaskCommand) NextStep() command.Command {
	return &closed_task.ClosedTaskByIDCommand{}
}

// InputField Поле для сохранения данных
func (c *ClosedTaskCommand) InputField() string {
	return "title"
}
