package closed_task

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
	btnclosedtaskbyid "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/closed_task"
	"github.com/pkg/errors"
)

// ClosedTaskByIdCommand - структура команды /closed_task_by_id
type ClosedTaskByIdCommand struct{}

func (c *ClosedTaskByIdCommand) Execute(s command.Model, msg messaging.Message) error {

	var (
		task *models.Task
	)

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := msg.UserDisplayName
	if len(displayName) == 0 {
		displayName = msg.UserName
	}

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTClosedTaskByIdCommand, displayName)

	session, err := s.GetSessionService().GetOrCreateSession(s.GetCtx(), msg.UserID)

	taskID, err := sessionutils.ExtractInt64FromSession(session, "payload_task_id")
	if err != nil {
		return err
	}

	task, err = s.GetTaskStorage().GetTaskByID(s.GetCtx(), taskID)
	if err != nil {
		return errors.Wrap(err, "Не удалось изменить статус задачи")
	}

	task.StatusID = repositories.STATUS_CLOSED
	task.UpdatedAt = time.Now()

	err = s.GetTaskStorage().UpdateTask(s.GetCtx(), task)
	if err != nil {
		return err
	}

	// Инициализируем кнопки
	var buttons []bottypes.TgRowButtons

	buttons = append(buttons, btnclosedtaskbyid.BtnClosedTaskByID...)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)

}
