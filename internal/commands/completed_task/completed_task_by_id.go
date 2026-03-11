package completed_task

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btncompletedtaskbyid "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/completedtask"
	"github.com/pkg/errors"
)

// CompletedTaskByIDCommand - структура команды /completed_task_by_id
// Команда для выполнения задачи по ID.
type CompletedTaskByIDCommand struct{}

// Execute — реализация команды /completed_task_by_id.
// Отправляет приветственное сообщение с inline-кнопками пользователю.
func (c *CompletedTaskByIDCommand) Execute(s command.Model, msg messaging.Message) error {

	var task *models.Task

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := helpers.GetDisplayName(msg)

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTCompletedTaskByIDCommand, displayName)

	taskID, err := strconv.ParseInt(msg.Arguments, 10, 64)
	if err != nil {
		return s.GetTgClient().SendMessage(resources.ErrInvalidTaskID, msg.UserID)
	}

	task, err = s.GetTaskStorage().GetTaskByID(s.GetCtx(), taskID)
	if err != nil {
		return errors.Wrap(err, resources.ErrFailedToUpdateTask)
	}

	task.StatusID = repositories.StatusCompleted
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

// NextStep Следующая команда
func (c *CompletedTaskByIDCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *CompletedTaskByIDCommand) InputField() string {
	return "task_id"
}
