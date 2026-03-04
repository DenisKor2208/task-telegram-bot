package edit_task

import (
	"fmt"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/callbacktokenpayloadutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/sessionutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	btnedittaskbyid "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/edit_task"
	"github.com/pkg/errors"
)

// EditTaskByIdCommand - структура команды /edit_task_by_id
type EditTaskByIdCommand struct{}

func (c *EditTaskByIdCommand) Execute(s command.Model, msg messaging.Message) error {

	var task *models.Task

	session, err := s.GetSessionService().GetOrCreateSession(s.GetCtx(), msg.UserID)

	args, err := sessionutils.ExtractInt64FromSession(session, "payload_task_id")
	if err != nil {
		return err
	}

	// Проверяем, что аргументы являются числом
	/* 	argsInt, err := strconv.Atoi(strings.TrimSpace(args))
	   	if err != nil {
	   		return errors.Wrap(err, "Не удалось удалить задачу")
	   	} */

	task, err = s.GetTaskStorage().GetTaskByID(s.GetCtx(), args)
	if err != nil {
		return errors.Wrap(err, "Не удалось изменить статус задачи")
	}
	text := fmt.Sprintf("%s %s", task.Description, task.Deadline.Format("02.01.2006 15:04"))

	// Инициализируем кнопки
	var buttons []bottypes.TgRowButtons

	// Создаём токен вместо строки "/edit_task_by_id 123"
	token, _ := s.GetSessionService().CreateCallbackToken(
		s.GetCtx(),
		callbacktokenpayloadutils.CallbackTokenPayload{
			Action: "update_task",
			TaskID: int64(task.ID),
			UserID: msg.UserID,
			Field:  "title",
		},
		10*time.Minute,
	)
	/* 	if err != nil {
		logger.Error("Failed to create token", "err", err)
	} */

	// Устанавливаем состояние сессии
	s.GetSessionService()
	sessionutils.SetSessionTaskToken(
		s.GetCtx(), s.GetSessionService(), msg.UserID, token)

	buttons = append(buttons, btnedittaskbyid.BtnEditTaskByID...)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)

}
