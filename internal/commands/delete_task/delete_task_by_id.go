package delete_task

import (
	"fmt"
	"strconv"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btndeletetaskbyid "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/delete_task"
)

// DeleteTaskByIDCommand - структура команды /delete_task_by_id
// Команда для удаления задачи по ID через callback.
type DeleteTaskByIDCommand struct{}

// Execute — реализация команды /delete_task_by_id.
// Отправляет приветственное сообщение с inline-кнопками пользователю.
func (c *DeleteTaskByIDCommand) Execute(s command.Model, msg messaging.Message) error {

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := helpers.GetDisplayName(msg)

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTDeleteTaskByIDCommand, displayName)

	taskID, err := strconv.ParseInt(msg.Arguments, 10, 64)
	if err != nil {
		return s.GetTgClient().SendMessage("Неверный ID задачи", msg.UserID)
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

// NextStep Следующая команда
func (c *DeleteTaskByIDCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *DeleteTaskByIDCommand) InputField() string {
	return "task_id"
}
