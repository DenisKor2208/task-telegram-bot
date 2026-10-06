package closed_task

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btnclosedtaskbyid "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/closedtask"
	"github.com/pkg/errors"
)

// ClosedTaskByIDCommand - структура команды /closed_task_by_id
type ClosedTaskByIDCommand struct{}

// Execute — реализация команды /closed_task_by_id.
func (c *ClosedTaskByIDCommand) Execute(s command.Model, msg messaging.Message) error {

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := helpers.GetDisplayName(msg)

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTClosedTaskByIDCommand, displayName)

	taskID, err := strconv.ParseInt(msg.Arguments, 10, 64)
	if err != nil {
		return s.GetTgClient().SendMessage(resources.ErrInvalidTaskID, msg.UserID)
	}

	task, err := s.GetTaskStorage().GetTaskByID(s.GetCtx(), msg.UserID, taskID)
	if errors.Is(err, repositories.ErrTaskNotFound) {
		// Задачи нет или она чужая — для пользователя это одно и то же
		return messaging.NewUserError(resources.ErrTaskNotFound)
	}
	if err != nil {
		return errors.Wrap(err, resources.ErrFailedToUpdateTask)
	}

	task.StatusID = repositories.StatusClosed
	task.UpdatedAt = time.Now()

	err = s.GetTaskStorage().UpdateTask(s.GetCtx(), task)
	if errors.Is(err, repositories.ErrTaskNotFound) {
		return messaging.NewUserError(resources.ErrTaskNotFound)
	}
	if err != nil {
		return err
	}

	// Инициализируем кнопки
	var buttons []bottypes.TgRowButtons

	buttons = append(buttons, btnclosedtaskbyid.BtnClosedTaskByID...)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)

}

// NextStep Следующая команда
func (c *ClosedTaskByIDCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *ClosedTaskByIDCommand) InputField() string {
	return "task_id"
}
