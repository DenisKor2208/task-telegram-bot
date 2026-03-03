package commands

import (
	"fmt"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/callbacktokenpayloadutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
	"github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands"

	"github.com/pkg/errors"
)

const (
	COMMAND_EDIT_TASKS = "/edit_task"
)

// EditTaskCommand - структура команды /edit_task - "Редактировать задачу"
type EditTaskCommand struct{}

func (c *EditTaskCommand) Execute(s types.Model, msg types.Message, session *models.UserSession) error {

	var tasks []*models.Task
	var err error

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := msg.UserDisplayName
	if len(displayName) == 0 {
		displayName = msg.UserName
	}

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTEditTask, displayName)

	statuses := []int{
		repositories.STATUS_IN_PROGRESS,
		repositories.STATUS_COMPLETED,
		repositories.STATUS_OVERDUE,
		repositories.STATUS_CLOSED,
	}

	tasks, err = s.GetTaskStorage().GetTasksByStatusID(s.GetCtx(), statuses)
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
		/*		if task.Deadline.IsZero() {
				deadlineStr = "Без дедлайна"
			}*/

		// Создаём токен вместо строки "/edit_task_by_id 123"
		token, err := s.GetSessionService().CreateCallbackToken(
			s.GetCtx(),
			callbacktokenpayloadutils.CallbackTokenPayload{
				Action: "edit_task_by_id",
				TaskID: int64(task.ID),
				UserID: msg.UserID,
				Field:  "title",
			},
			10*time.Minute,
		)
		if err != nil {
			logger.Error("Failed to create token", "err", err)
			continue // Пропускаем кнопку, если ошибка
		}

		// Создаём кнопку для задачи
		button := bottypes.TgInlineButton{
			DisplayName: fmt.Sprintf("%s (до %s)", task.Description, deadlineStr),
			Value:       fmt.Sprintf("/edit_task_by_id %s", token),
		}

		// Добавляем как новый ряд
		additionalButtons = append(additionalButtons, bottypes.TgRowButtons{button})
	}

	// Формируем общий список кнопок c дополнительными и базовыми кнопками
	buttons = append(buttons, additionalButtons...)
	buttons = append(buttons, commands.BtnEditTask...)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)

}
