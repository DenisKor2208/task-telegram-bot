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
func (c *ClosedTaskCommand) Execute(s command.Model, msg messaging.Message) error {

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := helpers.GetDisplayName(msg)

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTClosedTask, displayName)

	statuses := []int{
		repositories.StatusInProgress,
		repositories.StatusCompleted,
		repositories.StatusOverdue,
		// repositories.StatusClosed,
	}

	tasks, err := s.GetTaskStorage().GetTasksByStatusID(s.GetCtx(), msg.UserID, statuses)
	if err != nil {
		return errors.Wrap(err, "Не удалось получить задачи")
	}

	// Инициализируем кнопки
	var buttons []bottypes.TgRowButtons

	// Инициализируем пустой срез для дополнительных кнопок
	var additionalButtons []bottypes.TgRowButtons

	// Дедлайны показываем в часовом поясе пользователя
	loc := helpers.UserLocation(s.GetCtx(), s.GetUserStorage(), msg.UserID)

	// Добавляем кнопки для каждой задачи
	for _, task := range tasks {
		deadlineStr := helpers.DeadlineLabel(task.Deadline, loc)

		button := bottypes.TgInlineButton{
			DisplayName: fmt.Sprintf("%s (%s)", task.Description, deadlineStr),
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
