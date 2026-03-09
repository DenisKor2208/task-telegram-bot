package views_tasks

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btnviewstasks "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/views_tasks"
	"github.com/pkg/errors"
)

// FilterOverdueCommand - структура команды /filter_overdue
type FilterOverdueCommand struct{}

// Execute — реализация команды /filter_overdue.
// Отправляет приветственное сообщение с inline-кнопками пользователю.
func (c *FilterOverdueCommand) Execute(s command.Model, msg messaging.Message) error {

	statuses := []int{repositories.STATUS_OVERDUE}

	tasks, err := s.GetTaskStorage().GetTasksByStatusID(s.GetCtx(), statuses)
	if err != nil {
		return errors.Wrap(err, "Не удалось получить задачи")
	}

	var buttons []bottypes.TgRowButtons
	var taskButtons []bottypes.TgRowButtons

	// Цикл по задачам: добавляем кнопку для каждой (предполагаю, что TgRowButtons - это срез кнопок)
	for _, task := range tasks {
		deadlineStr := task.Deadline.Format("02.01.2006 15:04")
		if task.Deadline.IsZero() {
			deadlineStr = "Без дедлайна"
		}

		button := bottypes.TgInlineButton{
			DisplayName: fmt.Sprintf("%s (до %s)", task.Description, deadlineStr),
			Value:       fmt.Sprintf("/view_tasks %d", task.ID),
		}

		taskButtons = append(taskButtons, bottypes.TgRowButtons{button})
	}

	buttons = append(buttons, taskButtons...)
	buttons = append(buttons, btnviewstasks.BtnFilterCompleted...)

	// Добавляем кнопку "Назад"
	lastCommand := "view_tasks"
	buttons = append(
		buttons,
		bottypes.TgRowButtons{
			{DisplayName: actionBack, Value: "/" + lastCommand},
		})

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := helpers.GetDisplayName(msg)
	text := fmt.Sprintf(resources.TXTFilterOverdue, displayName)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)

}

// NextStep Следующая команда
func (c *FilterOverdueCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *FilterOverdueCommand) InputField() string {
	return ""
}
