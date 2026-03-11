// Package add_task
package add_task

import (
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	"github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands"
	"github.com/pkg/errors"
)

// SaveTaskCommand - структура команды /save_task
type SaveTaskCommand struct{}

// Execute — реализация команды /save_task.
// Отправляет приветственное сообщение с inline-кнопками пользователю.
func (c *SaveTaskCommand) Execute(s command.Model, msg messaging.Message) error {

	session, err := s.GetSessionService().GetOrCreateSession(s.GetCtx(), msg.UserID)
	if err != nil {
		return errors.Wrap(err, "не удалось сохранить задачу")
	}

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
	taskDesc, taskDate, ok := helpers.ParseForCommandSaveTask(argsText)
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
		Timezone:  "UTC",
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	// Сохраняем пользователя
	createdUser, err := s.GetUserStorage().CreateUser(s.GetCtx(), user)
	if err != nil {
		return errors.Wrap(err, "не удалось сохранить задачу")
	}

	// Загружаем часовой пояс пользователя
	loc, err := time.LoadLocation(createdUser.Timezone)
	if err != nil {
		logger.Warn("Не удалось загрузить часовой пояс, используется UTC",
			"user_id", createdUser.ID, "timezone", createdUser.Timezone, "error", err)
		loc = time.UTC
	}

	// Формируем модель задачи
	task := &models.Task{
		Description: taskDesc,
		UserID:      createdUser.ID,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}

	if taskDate != "" {
		layout := "02.01.2006 15:04"
		parsed, err := time.ParseInLocation(layout, taskDate, loc)
		if err != nil {
			logger.Warn("Не удалось распарсить дедлайн, будет проигнорирован", "date", taskDate, "error", err)
		} else {
			task.Deadline = parsed.UTC()
		}
	}

	// Определяем статус задачи на основе дедлайна
	nowUTC := time.Now().UTC()
	task.StatusID = repositories.StatusInProgress
	if task.Deadline.Before(nowUTC) {
		task.StatusID = repositories.StatusOverdue
	}

	_, err = s.GetTaskStorage().CreateTask(s.GetCtx(), task)
	if err != nil {
		return errors.Wrap(err, "Не удалось сохранить задачу")
	}

	return s.GetTgClient().ShowInlineButtons(resources.TXTSaveTask, commands.BtnSaveTask, msg.UserID)

}

// NextStep Следующая команда
func (c *SaveTaskCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *SaveTaskCommand) InputField() string {
	return ""
}
