package views_tasks

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btnviewstasks "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/views_tasks"
	"github.com/pkg/errors"
)

// FilterCompletedCommand - структура команды /filter_completed
type FilterCompletedCommand struct{}

func (c *FilterCompletedCommand) Execute(s command.Model, msg messaging.Message) error {

	var tasks []*models.Task
	var err error

	statuses := []int{
		repositories.STATUS_COMPLETED,
	}

	session, err := s.GetSessionService().GetOrCreateSession(s.GetCtx(), msg.UserID)

	tasks, err = s.GetTaskStorage().GetTasksByStatusID(s.GetCtx(), statuses)
	if err != nil {
		return errors.Wrap(err, "Не удалось получить задачи")
	}

	// Инициализируем кнопки с базовыми (BtnFilterCompleted)
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
	buttons = append(buttons, btnviewstasks.BtnFilterCompleted...)

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
	text := fmt.Sprintf(resources.TXTFilterCompleted, displayName)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)

}
