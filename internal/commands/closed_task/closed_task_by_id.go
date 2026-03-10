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
	btnclosedtaskbyid "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/closed_task"
	"github.com/pkg/errors"
)

// ClosedTaskByIDCommand - структура команды /closed_task_by_id
// Команда для завершения задачи по ID.
type ClosedTaskByIDCommand struct{}

// Execute — реализация команды /closed_task_by_id.
// Отправляет приветственное сообщение с inline-кнопками пользователю.
func (c *ClosedTaskByIDCommand) Execute(s command.Model, msg messaging.Message) error {

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := helpers.GetDisplayName(msg)

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTClosedTaskByIDCommand, displayName)

	taskID, err := strconv.ParseInt(msg.Arguments, 10, 64)
	if err != nil {
		return s.GetTgClient().SendMessage("Неверный ID задачи", msg.UserID)
	}

	task, err := s.GetTaskStorage().GetTaskByID(s.GetCtx(), taskID)
	if err != nil {
		return errors.Wrap(err, "Не удалось изменить статус задачи")
	}

	task.StatusID = repositories.StatusClosed
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

// NextStep Следующая команда
func (c *ClosedTaskByIDCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *ClosedTaskByIDCommand) InputField() string {
	return "task_id"
}
