package delete_task

import (
	"fmt"
	"html"
	"strconv"

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

	taskID, err := strconv.ParseInt(msg.Arguments, 10, 64)
	if err != nil {
		return s.GetTgClient().SendMessage(resources.ErrInvalidTaskID, msg.UserID)
	}

	// Загружаем задачу до удаления
	task, err := s.GetTaskStorage().GetTaskByID(s.GetCtx(), msg.UserID, taskID)
	if errors.Is(err, repositories.ErrTaskNotFound) {
		// Задачи нет или она чужая — для пользователя это одно и то же
		return messaging.NewUserError(resources.ErrTaskNotFound)
	}
	if err != nil {
		return errors.Wrap(err, resources.ErrFailedToDeleteTask)
	}

	err = s.GetTaskStorage().DeleteTaskByID(s.GetCtx(), msg.UserID, taskID)
	if errors.Is(err, repositories.ErrTaskNotFound) {
		// Задачи нет или она чужая — для пользователя это одно и то же
		return messaging.NewUserError(resources.ErrTaskNotFound)
	}
	if err != nil {
		return err
	}

	// Описание экранируется
	text := fmt.Sprintf(resources.TXTDeleteTaskByIDCommand, html.EscapeString(task.Description))

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
