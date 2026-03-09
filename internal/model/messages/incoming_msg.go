// Package messages
package messages

import (
	"context"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/sessionutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/service"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/storage"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// Константы для ключей сессии
const (
	SessionKeyLastCommand = "last_command"
	SessionKeyNextCommand = "next_command"
	SessionKeyTaskTitle   = "title"
)

// Model Модель бота (клиент, хранилище, последние команды пользователя)
type Model struct {
	ctx             context.Context
	userStorage     storage.User
	statusStorage   storage.Status
	taskStorage     storage.Task
	tgClient        messaging.Sender
	commandRegistry command.Registry
	sessionService  service.Session
}

// New Генерация сущности для хранения клиента ТГ и хранилища пользователей.
func New(
	ctx context.Context,
	userStorage storage.User,
	statusStorage storage.Status,
	taskStorage storage.Task,
	tgClient messaging.Sender,
	registry command.Registry,
	sessionService service.Session,
) *Model {
	return &Model{
		ctx:             ctx,
		userStorage:     userStorage,
		statusStorage:   statusStorage,
		taskStorage:     taskStorage,
		tgClient:        tgClient,
		commandRegistry: registry,
		sessionService:  sessionService,
	}
}

func (s *Model) GetCtx() context.Context              { return s.ctx }
func (s *Model) SetCtx(ctx context.Context)           { s.ctx = ctx }
func (s *Model) GetTgClient() messaging.Sender        { return s.tgClient }
func (s *Model) GetCommandRegistry() command.Registry { return s.commandRegistry }
func (s *Model) GetUserStorage() storage.User         { return s.userStorage }
func (s *Model) GetStatusStorage() storage.Status     { return s.statusStorage }
func (s *Model) GetTaskStorage() storage.Task         { return s.taskStorage }
func (s *Model) GetSessionService() service.Session   { return s.sessionService }

// IncomingMessage — главная точка входа для обработки сообщений от Telegram.
func (s *Model) IncomingMessage(msg messaging.Message) error {
	if msg.IsCallback || msg.IsCommand {
		return s.handleCallback(msg)
	}
	return s.handleUserInput(msg)
}

// handleCallback обрабатывает нажатия на inline-кнопки.
func (s *Model) handleCallback(msg messaging.Message) error {
	cleanSession := func() (*models.UserSession, error) {
		return sessionutils.ClearSessionData(s.ctx, s.sessionService, msg.UserID)
	}

	// Получаем команду
	cmdName := msg.Command
	if cmdName == "" {
		// Если не удалось распарсить команду, игнорируем (такого быть не должно)
		logger.Warn("Callback without command", "text", msg.Text)

		if _, err := cleanSession(); err != nil {
			logger.Error("Failed to clean session", "error", err)
		}
	}

	cmd, exists := s.commandRegistry.GetCommand(cmdName)
	if !exists {
		logger.Warn("Callback command not found", "command", cmdName)

		if _, err := cleanSession(); err != nil {
			logger.Error("Failed to clean session", "error", err)
		}

		return s.tgClient.SendMessage(resources.TXTUnknownCommand, msg.UserID)
	}

	_, saveErr := sessionutils.UpdateSessionData(
		s.ctx,
		s.sessionService,
		msg.UserID,
		map[string]any{SessionKeyLastCommand: cmdName})
	if saveErr != nil {
		logger.Error("Failed to save last command to session", "error", saveErr)
	}

	return cmd.Execute(s, msg)
}

// handleUserInput обрабатывает обычный текст (не команду) — это ввод для интерактивных команд.
func (s *Model) handleUserInput(msg messaging.Message) error {
	cleanSession := func() (*models.UserSession, error) {
		return sessionutils.ClearSessionData(s.ctx, s.sessionService, msg.UserID)
	}

	session, err := s.sessionService.GetOrCreateSession(s.ctx, msg.UserID)
	if err != nil {
		return err
	}

	lastCmdName, ok := session.Data[SessionKeyLastCommand].(string)
	if !ok {
		// Нет ожидающей команды — вероятно, пользователь просто написал текст без команды
		return s.tgClient.SendMessage("Я не понимаю. Введите /start для начала работы.", msg.UserID)
	}

	cmd, exists := s.commandRegistry.GetCommand(lastCmdName)
	if !exists {
		// Команда из сессии не найдена — очищаем сессию
		if _, err := cleanSession(); err != nil {
			logger.Error("Failed to clean session", "error", err)
		}

		return s.tgClient.SendMessage("Сессия устарела. Начните заново.", msg.UserID)
	}

	interactiveCmd, isInteractive := cmd.(command.InteractiveCommand)
	if !isInteractive {
		// Последняя команда не была интерактивной — очищаем сессию
		if _, err := cleanSession(); err != nil {
			logger.Error("Failed to clean session", "error", err)
		}

		return s.tgClient.SendMessage("Неожиданный ввод. Пожалуйста, введите команду.", msg.UserID)
	}

	// Сохраняем введённые данные под ключом, который возвращает InputField
	field := interactiveCmd.InputField()
	_, err = sessionutils.UpdateSessionData(s.ctx, s.sessionService, msg.UserID, map[string]any{field: msg.Text})
	if err != nil {
		logger.Error("Failed to save user input to session", "error", err)

		return s.tgClient.SendMessage("Ошибка при сохранении данных. Попробуйте позже.", msg.UserID)
	}

	// Переходим к следующему шагу
	nextCmd := interactiveCmd.NextStep()
	if nextCmd == nil {
		// Интерактив завершён — очищаем сессию и все данные сессии, связанные с этим интерактивом
		if _, err := cleanSession(); err != nil {
			logger.Error("Failed to clean session", "error", err)
		}
		return s.tgClient.SendMessage("Готово!", msg.UserID)
	}

	// Выполняем следующую команду, передавая то же сообщение
	return nextCmd.Execute(s, msg)
}
