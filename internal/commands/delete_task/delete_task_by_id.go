package delete_task

import (
	"fmt"
	"strconv"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btndeletetaskbyid "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/deletetask"
	"github.com/pkg/errors"
)

// DeleteTaskByIDCommand - структура команды /delete_task_by_id
type DeleteTaskByIDCommand struct{}

// Execute — реализация команды /delete_task_by_id.
func (c *DeleteTaskByIDCommand) Execute(s command.Model, msg messaging.Message) error {

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := helpers.GetDisplayName(msg)

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTDeleteTaskByIDCommand, displayName)

	taskID, err := strconv.ParseInt(msg.Arguments, 10, 64)
	if err != nil {
		return s.GetTgClient().SendMessage(resources.ErrInvalidTaskID, msg.UserID)
	}

	err = s.GetTaskStorage().DeleteTaskByID(s.GetCtx(), msg.UserID, taskID)
	if errors.Is(err, repositories.ErrTaskNotFound) {
		// Задачи нет или она чужая — для пользователя это одно и то же
		return messaging.NewUserError(resources.ErrTaskNotFound)
	}
	if err != nil {
		return err
	}

	// Инициализируем кнопки
	var buttons []bottypes.TgRowButtons

	buttons = append(buttons, btndeletetaskbyid.BtnDeleteTaskByID...)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)

}

// NextStep Следующая команда
func (c *DeleteTaskByIDCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *DeleteTaskByIDCommand) InputField() string {
	return "task_id"
}
