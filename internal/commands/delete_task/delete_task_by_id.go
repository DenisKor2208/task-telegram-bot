package commands

import (
	"fmt"
	"strconv"

	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
	btndeletetaskbyid "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/delete_task"
	"github.com/pkg/errors"
)

// DeleteTaskByIdCommand - структура команды /delete_task_by_id
type DeleteTaskByIdCommand struct{}

func (c *DeleteTaskByIdCommand) Execute(s types.Model, msg types.Message, session *models.UserSession) error {

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := msg.UserDisplayName
	if len(displayName) == 0 {
		displayName = msg.UserName
	}

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTDeleteTaskByIdCommand, displayName)

	// Парсим команду и аргументы
	taskId, ok := session.Data["payload_task_id"].(int64)
	if !ok {
		if legacyID, ok := session.Data["task_id"].(string); ok {
			id, err := strconv.ParseInt(legacyID, 10, 64)
			if err != nil {
				return errors.New("Не удалось удалить задачу")
			}
			taskId = id
		} else {
			return errors.New("ID задачи не найден в сессии")
		}
	}

	err := s.GetTaskStorage().DeleteTaskByID(s.GetCtx(), taskId)
	if err != nil {
		return err
	}

	// Инициализируем кнопки
	var buttons []bottypes.TgRowButtons

	buttons = append(buttons, btndeletetaskbyid.BtnDeleteTaskByID...)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)

}
