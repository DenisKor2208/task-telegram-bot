package commands

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/sessionutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btndeletetaskbyid "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/delete_task"
)

// DeleteTaskByIdCommand - структура команды /delete_task_by_id
type DeleteTaskByIdCommand struct{}

func (c *DeleteTaskByIdCommand) Execute(s command.Model, msg messaging.Message) error {

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := msg.UserDisplayName
	if len(displayName) == 0 {
		displayName = msg.UserName
	}

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTDeleteTaskByIdCommand, displayName)

	session, err := s.GetSessionService().GetOrCreateSession(s.GetCtx(), msg.UserID)

	taskID, err := sessionutils.ExtractInt64FromSession(session, "payload_task_id")
	if err != nil {
		return err
	}

	err = s.GetTaskStorage().DeleteTaskByID(s.GetCtx(), taskID)
	if err != nil {
		return err
	}

	// Инициализируем кнопки
	var buttons []bottypes.TgRowButtons

	buttons = append(buttons, btndeletetaskbyid.BtnDeleteTaskByID...)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)

}
