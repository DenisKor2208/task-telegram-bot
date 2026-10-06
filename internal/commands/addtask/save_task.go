// Package addtask
package addtask

import (
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/sessionutils"
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
func (c *SaveTaskCommand) Execute(s command.Model, msg messaging.Message) error {

	session, err := s.GetSessionService().GetOrCreateSession(s.GetCtx(), msg.UserID)
	if err != nil {
		return errors.Wrap(err, resources.ErrFailedToSaveTask)
	}

	// Проверяем наличие аргументов только в сессии
	var argsText string
	if session.Data == nil {
		return errors.New("no data in session: arguments not found")
	}
	if args, ok := session.Data["title"].(string); ok && args != "" {
		argsText = args
	} else {
		// Пользователь прислал не текст (фото, стикер и т.п.) — он может исправить это сам
		return messaging.NewUserError(resources.ErrTaskDescriptionEmpty)
	}

	// Теперь парсим аргументы из сессии
	input, ok := helpers.ParseTaskInput(argsText)
	if !ok {
		return messaging.NewUserError(resources.ErrTaskDescriptionEmpty)
	}

	// В конце указано время, но без даты (например, «Купить молоко в 18:00»)
	if input.TimeWithoutDate {
		return messaging.NewUserError(resources.ErrTimeWithoutDate)
	}

	// Если описание задачи пустое (например, введена только дата)
	if input.Description == "" {
		return messaging.NewUserError(resources.ErrTaskDescriptionEmpty)
	}

	// Формируем модель пользователя
	timestamp := time.Unix(msg.Date, 0)
	user := &models.User{
		TgID:      int(msg.UserID),
		Name:      helpers.GetUserNameForDB(msg),
		Timezone:  "UTC",
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	// Получаем пользователя (обычно он уже создан в Model.IncomingMessage;
	// CreateUser вернёт существующего или создаст, если там не получилось)
	createdUser, err := s.GetUserStorage().CreateUser(s.GetCtx(), user)
	if err != nil {
		return errors.Wrap(err, resources.ErrFailedToSaveTask)
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
		Description: input.Description,
		UserID:      createdUser.ID,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}

	if input.Date != "" {
		// Указана только дата — дедлайн до конца дня
		timeStr := input.Time
		if timeStr == "" {
			timeStr = helpers.DefaultDeadlineTime
		}

		// Дата и время интерпретируются в часовом поясе пользователя
		parsed, err := time.ParseInLocation(helpers.DeadlineInputLayout, input.Date+" "+timeStr, loc)
		if err != nil {
			// Формат верный, но такой даты или времени нет (например, 31.02.2026 или 25:00)
			return messaging.NewUserError(resources.ErrInvalidDeadline)
		}

		if !parsed.After(time.Now()) {
			return messaging.NewUserError(resources.ErrDeadlinePassed)
		}

		deadlineUTC := parsed.UTC()
		task.Deadline = &deadlineUTC
	}

	// Определяем статус задачи на основе дедлайна (задача без дедлайна не может быть просрочена)
	nowUTC := time.Now().UTC()
	task.StatusID = repositories.StatusInProgress
	if task.Deadline != nil && task.Deadline.Before(nowUTC) {
		task.StatusID = repositories.StatusOverdue
	}

	_, err = s.GetTaskStorage().CreateTask(s.GetCtx(), task)
	if err != nil {
		return errors.Wrap(err, resources.ErrFailedToSaveTask)
	}

	// Задача сохранена — диалог добавления завершён. Очищаем сессию,
	// чтобы следующий текст пользователя не создал ещё одну задачу.
	if _, err := sessionutils.ClearSessionData(s.GetCtx(), s.GetSessionService(), msg.UserID); err != nil {
		logger.Error("Не удалось очистить сессию после сохранения задачи", "user_id", msg.UserID, "error", err)
	}

	// Показываем, как бот понял дедлайн, чтобы пользователь сразу заметил ошибку
	text := resources.TXTSaveTask + "\n" + helpers.DeadlineText(task.Deadline, loc)

	return s.GetTgClient().ShowInlineButtons(text, commands.BtnSaveTask, msg.UserID)

}

// NextStep Следующая команда
func (c *SaveTaskCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *SaveTaskCommand) InputField() string {
	return ""
}
