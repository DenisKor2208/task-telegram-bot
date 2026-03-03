package views_tasks

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
	btnviewstasks "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/views_tasks"
	"github.com/pkg/errors"
)

// FilterInProgressCommand - структура команды /filter_in_progress
type FilterInProgressCommand struct{}

func (c *FilterInProgressCommand) Execute(s types.Model, msg types.Message, session *models.UserSession) error {

	var tasks []*models.Task
	var err error

	statuses := []int{
		repositories.STATUS_IN_PROGRESS,
	}

	tasks, err = s.GetTaskStorage().GetTasksByStatusID(s.GetCtx(), statuses)
	if err != nil {
		return errors.Wrap(err, "Не удалось получить задачи")
	}

	// Инициализируем кнопки с базовыми (BtnFilterInProgress)
	var buttons []bottypes.TgRowButtons

	// Инициализируем пустой срез для дополнительных кнопок (один из вариантов: var)
	var additionalButtons []bottypes.TgRowButtons

	// Цикл по задачам: добавляем кнопку для каждой (предполагаю, что TgRowButtons - это срез кнопок)
	for _, task := range tasks {
		deadlineStr := task.Deadline.Format("02.01.2006 15:04")
		/*		if task.Deadline.IsZero() {
				deadlineStr = "Без дедлайна"
			}*/

		// Создаём кнопку для задачи (DisplayName - имя + дедлайн, Value - команда с ID)
		button := bottypes.TgInlineButton{
			DisplayName: fmt.Sprintf("%s (до %s)", task.Description, deadlineStr),
			Value:       fmt.Sprintf("/task %d", task.ID),
		}

		// Добавляем как новый ряд (TgRowButtons) - предполагаю, что каждый ряд - одна кнопка
		additionalButtons = append(additionalButtons, bottypes.TgRowButtons{button})
	}

	// Формируем общий список кнопок с базовыми и дополнительными кнопками
	buttons = append(buttons, additionalButtons...)
	buttons = append(buttons, btnviewstasks.BtnFilterInProgress...)

	// Добавляем кнопку "Назад"
	lastCommand := session.Data["last_command"].(string)
	buttons = append(
		buttons,
		bottypes.TgRowButtons{
			{DisplayName: "Назад", Value: "/" + lastCommand},
		})

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := msg.UserDisplayName
	if len(displayName) == 0 {
		displayName = msg.UserName
	}

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTFilterInProgress, displayName)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)

}
