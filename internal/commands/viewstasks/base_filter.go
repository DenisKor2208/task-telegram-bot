// Package viewstasks
package viewstasks

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	btnbasefilter "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/viewstasks"
	"github.com/pkg/errors"
)

// BaseFilterCommand - общая команда для фильтрации задач по статусам.
type BaseFilterCommand struct {
	Statuses  []int  // список статусов
	Text      string // заголовок списка
	EmptyText string // текст, если задач с такими статусами нет
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

	// Дедлайны показываем в часовом поясе пользователя
	loc := helpers.UserLocation(s.GetCtx(), s.GetUserStorage(), msg.UserID)

	for _, task := range tasks {
		button := helpers.TaskButton(task, loc, "/view_tasks")
		taskButtons = append(taskButtons, bottypes.TgRowButtons{button})
	}

	buttons = append(buttons, taskButtons...)
	buttons = append(buttons, btnbasefilter.BtnBaseFilter...)

	displayName := helpers.GetDisplayName(msg)
	text := fmt.Sprintf(c.Text, displayName)
	if len(tasks) == 0 {
		text = fmt.Sprintf(c.EmptyText, displayName)
	}

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
