package completed_task

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
	btncompletedtaskbyid "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/completed_task"
	"github.com/pkg/errors"
)

// CompletedTaskByIdCommand - структура команды /completed_task_by_id
type CompletedTaskByIdCommand struct{}

func (c *CompletedTaskByIdCommand) Execute(s types.Model, msg types.Message, session *models.UserSession) error {

	var task *models.Task

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := msg.UserDisplayName
	if len(displayName) == 0 {
		displayName = msg.UserName
	}

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTCompletedTaskByIdCommand, displayName)
	taskId, ok := session.Data["payload_task_id"].(int64)
	if !ok {
		// Fallback
		if legacyID, ok := session.Data["task_id"].(string); ok {
			id, err := strconv.ParseInt(legacyID, 10, 64)
			if err != nil {
				return errors.New("Не удалось выполнить задачу")
			}
			taskId = id
		} else {
			return errors.New("ID задачи не найден в сессии")
		}
	}

	task, err := s.GetTaskStorage().GetTaskByID(s.GetCtx(), taskId)
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
