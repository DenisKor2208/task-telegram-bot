package views_tasks

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btnviewstasks "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/views_tasks"
	"github.com/pkg/errors"
)

const (
	actionBack = "Назад"
)

// FilterAllCommand - структура команды /filter_all
type FilterAllCommand struct{}

// Execute — реализация команды /filter_all.
// Отправляет приветственное сообщение с inline-кнопками пользователю.
func (c *FilterAllCommand) Execute(s command.Model, msg messaging.Message) error {

	tasks, err := s.GetTaskStorage().GetAllTasks(s.GetCtx())
	if err != nil {
		return errors.Wrap(err, "Не удалось получить задачи")
	}

	var buttons []bottypes.TgRowButtons
	var taskButtons []bottypes.TgRowButtons

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
	buttons = append(buttons, btnviewstasks.BtnFilterAll...)

	// Добавляем кнопку "Назад"
	lastCommand := "view_tasks"

	backBtn := bottypes.TgRowButtons{
		{DisplayName: actionBack, Value: "/" + lastCommand},
	}
	buttons = append(buttons, backBtn)

	displayName := helpers.GetDisplayName(msg)
	text := fmt.Sprintf(resources.TXTFilterAll, displayName)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)
}

// NextStep Следующая команда
func (c *FilterAllCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *FilterAllCommand) InputField() string {
	return ""
}
