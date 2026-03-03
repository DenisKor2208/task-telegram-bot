package add_task

import (
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/model/messages"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
	"github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands"
	"github.com/pkg/errors"
)

// SaveTaskCommand - структура команды /save_task
type SaveTaskCommand struct{}

func (c *SaveTaskCommand) Execute(s types.Model, msg types.Message, session *models.UserSession) error {

	// Проверяем наличие аргументов только в сессии
	var argsText string
	if session.Data == nil {
		return errors.New("no data in session: arguments not found")
	}
	if args, ok := session.Data["title"].(string); ok && args != "" {
		argsText = args
	} else {
		return errors.New("arguments not found or invalid in session")
	}

	// Теперь парсим аргументы из сессии
	taskDesc, taskDate, ok := messages.ParseForCommandSaveTask(argsText)
	if !ok {
		return errors.New("Не удалось сохранить задачу")
	}

	// Если описание задачи пустое
	if taskDesc == "" {
		return errors.New("Описание задачи отсутствует")
	}

	// Формируем модель пользователя
	timestamp := time.Unix(msg.Date, 0)
	user := &models.User{
		TgID:      int(msg.UserID),
		Name:      msg.UserName,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	// Сохраняем пользователя
	createdUser, err := s.GetUserStorage().CreateUser(s.GetCtx(), user)
	if err != nil {
		return errors.Wrap(err, "не удалось сохранить задачу")
	}

	// Назначаем задаче статус "В процессе"
	statusInProgress := repositories.STATUS_IN_PROGRESS

	// Формируем модель задачи
	task := &models.Task{
		Description: taskDesc,
		StatusID:    statusInProgress,
		UserID:      createdUser.ID,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}

	if taskDate != "" {
		layout := "02.01.2006 15:04"
		task.Deadline, _ = time.Parse(layout, taskDate)
	}

	_, err = s.GetTaskStorage().CreateTask(s.GetCtx(), task)
	if err != nil {
		return errors.Wrap(err, "Не удалось сохранить задачу")
	}

	// _ = s.GetSessionService().DeleteSession(s.GetCtx(), msg.UserID)
	//if err != nil {
	//	return errors.Wrap(err, "Не удалось сохранить задачу")
	//}

	return s.GetTgClient().ShowInlineButtons(resources.TXTSaveTask, commands.BtnSaveTask, msg.UserID)

}
