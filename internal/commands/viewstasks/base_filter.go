// Package viewstasks
package viewstasks

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btnbasefilter "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/viewstasks"
	"github.com/pkg/errors"
)

// BaseFilterCommand - общая команда для фильтрации задач по статусам.
type BaseFilterCommand struct {
	Statuses    []int  // список статусов
	ResourceKey string // ключ для получения текста из resources.FilterTexts
}

func (c *BaseFilterCommand) Execute(s command.Model, msg messaging.Message) error {
	var tasks []*models.Task
	var err error

	if len(c.Statuses) == 0 {
		tasks, err = s.GetTaskStorage().GetAllTasks(s.GetCtx(), msg.UserID)
	} else {
		tasks, err = s.GetTaskStorage().GetTasksByStatusID(s.GetCtx(), msg.UserID, c.Statuses)
	}
	if err != nil {
		return errors.Wrap(err, "Не удалось получить задачи")
	}

	var buttons []bottypes.TgRowButtons
	var taskButtons []bottypes.TgRowButtons

	for _, task := range tasks {
		deadlineStr := "Без дедлайна"
		if task.Deadline != nil {
			deadlineStr = task.Deadline.Format("02.01.2006 15:04")
		}

		button := bottypes.TgInlineButton{
			DisplayName: fmt.Sprintf("%s (до %s)", task.Description, deadlineStr),
			Value:       fmt.Sprintf("/view_tasks %d", task.ID),
		}
		taskButtons = append(taskButtons, bottypes.TgRowButtons{button})
	}

	buttons = append(buttons, taskButtons...)
	buttons = append(buttons, btnbasefilter.BtnBaseFilter...)

	displayName := helpers.GetDisplayName(msg)
	text := fmt.Sprintf(resources.TXTFilterAll, displayName)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)
}

// NextStep Следующая команда
func (c *BaseFilterCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *BaseFilterCommand) InputField() string {
	return ""
}
