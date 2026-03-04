package completed_task

import (
	"fmt"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/sessionutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btncompletedtaskbyid "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/completed_task"
	"github.com/pkg/errors"
)

// CompletedTaskByIdCommand - структура команды /completed_task_by_id
type CompletedTaskByIdCommand struct{}

func (c *CompletedTaskByIdCommand) Execute(s command.Model, msg messaging.Message) error {

	var task *models.Task

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := msg.UserDisplayName
	if len(displayName) == 0 {
		displayName = msg.UserName
	}

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTCompletedTaskByIdCommand, displayName)

	session, err := s.GetSessionService().GetOrCreateSession(s.GetCtx(), msg.UserID)

	taskID, err := sessionutils.ExtractInt64FromSession(session, "payload_task_id")
	if err != nil {
		return err
	}

	task, err = s.GetTaskStorage().GetTaskByID(s.GetCtx(), taskID)
	if err != nil {
		return errors.Wrap(err, "Не удалось изменить статус задачи")
	}

	task.StatusID = repositories.STATUS_COMPLETED
	task.UpdatedAt = time.Now()

	err = s.GetTaskStorage().UpdateTask(s.GetCtx(), task)
	if err != nil {
		return err
	}

	// Инициализируем кнопки
	var buttons []bottypes.TgRowButtons

	buttons = append(buttons, btncompletedtaskbyid.BtnCompletedTaskByID...)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)

}
