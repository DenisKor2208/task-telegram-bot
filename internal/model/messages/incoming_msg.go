package messages

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	configCommands "github.com/DenisKor2208/task-telegram-bot/internal/bot/config"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/sessionutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/service"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/storage"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	"github.com/pkg/errors"
)

// Model Модель бота (клиент, хранилище, последние команды пользователя)
type Model struct {
	ctx                context.Context
	userStorage        storage.User
	statusStorage      storage.Status
	taskStorage        storage.Task
	tgClient           messaging.Sender // Клиент
	commandRegistry    command.Registry // Хранилище команд
	lastUserCommand    map[int64]string // Последняя выбранная пользователем команда
	sessionService     service.Session
	configEntryService service.ConfigEntryService
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
		lastUserCommand: map[int64]string{},
		sessionService:  sessionService,
	}
}

func (s *Model) GetCtx() context.Context {
	return s.ctx
}

func (s *Model) SetCtx(ctx context.Context) {
	s.ctx = ctx
}

func (s *Model) GetTgClient() messaging.Sender {
	return s.tgClient
}

func (s *Model) GetCommandRegistry() command.Registry {
	return s.commandRegistry
}

func (s *Model) GetUserStorage() storage.User {
	return s.userStorage
}

func (s *Model) GetStatusStorage() storage.Status {
	return s.statusStorage
}

func (s *Model) GetTaskStorage() storage.Task {
	return s.taskStorage
}

func (s *Model) GetLastUserCommand(userID int64) string {
	return s.lastUserCommand[userID]
}

func (s *Model) SetLastUserCommand(userID int64, command string) {
	s.lastUserCommand[userID] = command
}

func (s *Model) GetSessionService() service.Session {
	return s.sessionService
}

func (s *Model) GetConfigEntryService() service.ConfigEntryService {
	return s.configEntryService
}

func (s *Model) SetConfigEntryService(configEntryService *configCommands.ConfigEntryService) {
	s.configEntryService = configEntryService
}

// IncomingMessage Обработка входящего сообщения.
func (s *Model) IncomingMessage(msg messaging.Message) error {

	// Распознавание стандартных команд.
	if isNeedReturn, err := checkBotCommands(s, msg); err != nil || isNeedReturn {
		return err
	}

	// Отправка ответа по умолчанию.
	return s.tgClient.SendMessage(resources.TXTUnknownCommand, msg.UserID)
}

// Распознавание стандартных команд бота.
func checkBotCommands(s command.Model, msg messaging.Message) (bool, error) {

	// Если команда или аргументы не указаны, то не обрабатываем.
	if msg.Command == "" && msg.Arguments == "" {
		return false, errors.New("message is not a command or arguments")
	}

	// Получаем сессию пользователя
	session, err := s.GetSessionService().GetOrCreateSession(s.GetCtx(), msg.UserID)
	if err != nil {
		logger.Error("Error getting user session", "err", err, "userID", msg.UserID)
		return false, err
	}
	// Если не команда, без аргумента и в сессии нету токена задачи
	if !msg.IsCommand && msg.Arguments != "" && session.TaskToken != "" {
		msg.Arguments = session.TaskToken
	}

	// Проверяем, не является ли аргумент токеном задачи
	if s.GetSessionService().IsCallbackToken(msg.Arguments) {
		return handleCallbackToken(s, msg)
	}

	// Обработка команд
	if msg.IsCommand {
		return handleCommand(s, msg)
	}

	// Обработка
	return handleFSM(s, msg, session)
}

// ParseCommandAndArgs парсит входную строку и извлекает команду (начинающуюся с "/") и аргументы.
// Возвращает команду и аргументы
func ParseCommandAndArgs(input string) (command string, args string) {
	trimmedInput := strings.TrimSpace(input)

	if trimmedInput == "" {
		return "", ""
	}

	//re := regexp.MustCompile(`^/([^\s]+)(?:\s+(.*))?`) // старый
	re := regexp.MustCompile(`^/([^\s]+)\s*(.*)`) // новый
	matches := re.FindStringSubmatch(trimmedInput)

	if len(matches) > 0 {
		command = "/" + strings.TrimSpace(matches[1])
		args = strings.TrimSpace(matches[2])
	} else {
		command = ""
		args = trimmedInput
	}

	return command, args
}

// ParseForCommandSaveTask парсит строку в формате "/save_task <описание> <дата в DD.MM.YYYY> <время в HH:MM>".
// Возвращает описание, дату, время и флаг успешности (true, если парсинг удался).
func ParseForCommandSaveTask(input string) (description, date string, ok bool) {

	pattern := `^(.+)\s+(\d{2}\.\d{2}\.\d{4})\s+(\d{2}:\d{2})$`
	re := regexp.MustCompile(pattern)

	// Убираем лишние пробелы по краям
	trimmed := strings.TrimSpace(input)
	matches := re.FindStringSubmatch(trimmed)

	if len(matches) != 4 {
		return "", "", false
	}

	description = strings.TrimSpace(matches[1])
	datetime := matches[2] + " " + matches[3]

	return description, datetime, true
}

// handleCommand — обрабатывает обычные команды (/start, /edit_task, /delete_task...)
func handleCommand(s command.Model, msg messaging.Message) (bool, error) {

	var err error

	for cmdName, config := range s.GetConfigEntryService().GetConfig() {
		if msg.Command == cmdName {
			if cmd, exists := s.GetCommandRegistry().GetCommand(cmdName); exists {

				nextState := config.TargetCommand

				// Устанавливаем в сессию потенциальную следующую команду
				if nextState != "" {
					_, _ = sessionutils.UpdateSessionData(s.GetCtx(), s.GetSessionService(), msg.UserID, map[string]any{"next_command": nextState})
				}

				err = cmd.Execute(s, msg)
				if err != nil {
					return false, err
				}

				// Сохраняем текущую выполненную команду
				_, _ = sessionutils.UpdateSessionData(s.GetCtx(), s.GetSessionService(), msg.UserID, map[string]any{"last_command": cmdName})

				return true, nil
			}
		}
	}

	// Команда не найдена
	logger.Warn("Unknown command", "command", msg.Command)

	return false, nil
}

// handleFSM — обрабатывает аргументы без команды (продолжение диалога)
// Например: после /add_task пользователь отправляет текст задачи
func handleFSM(s command.Model, msg messaging.Message, session *models.UserSession) (bool, error) {

	// Предыдущая выполненная команда
	lastCommand, _ := session.Data["last_command"].(string)
	if lastCommand == "" {
		return false, nil
	}

	// Следующая предполагаемая команда для выполнения
	nextCommand, _ := session.Data["next_command"].(string)
	if nextCommand == "" {
		return false, nil
	}

	// Конфиги всех комманд
	commandConfig := s.GetConfigEntryService().GetConfig()

	// Конфиг предполагаемая команда для выполнения
	config, ok := commandConfig[lastCommand]
	if !ok {
		// Неизвестное состояние целевая команда
		logger.Warn("Unknown target command", "target command", lastCommand)
		return false, nil
	}

	// Проверяем, есть ли команда в реестре
	cmd, exists := s.GetCommandRegistry().GetCommand(nextCommand)
	if !exists {
		logger.Error("Target command not found", "command", nextCommand)
		return false, errors.New("command not found")
	}

	var err error

	// Сохраняем аргументы в сессию (в поле, указанное в конфиге) //TODO Скорее всего нужно перебирать в цикле
	_, _ = sessionutils.UpdateSessionData(
		s.GetCtx(), s.GetSessionService(), msg.UserID,
		map[string]any{config.CommandFields[0]: msg.Arguments},
	)

	// Выполняем целевую команду
	err = cmd.Execute(s, msg)
	if err != nil {
		return false, err
	}

	_, _ = sessionutils.SetSessionTaskToken(s.GetCtx(), s.GetSessionService(), msg.UserID, "")
	_, _ = sessionutils.ClearSessionData(s.GetCtx(), s.GetSessionService(), msg.UserID)
	_, _ = sessionutils.UpdateSessionData(s.GetCtx(), s.GetSessionService(), msg.UserID, map[string]any{"last_command": nextCommand})

	// Некоторых команд удаляем сессию полностью (диалог окончен)
	/* 	if slices.Contains([]string{"save_task", "update_task"}, config.TargetCommand) {
		_ = s.GetSessionService().DeleteSession(s.GetCtx(), msg.UserID)
	} */

	return true, nil
}

// handleCallbackToken — обрабатывает нажатие кнопки с callback-токеном
func handleCallbackToken(s command.Model, msg messaging.Message) (bool, error) {
	// Получаем данные задачи по токену
	payload, err := s.GetSessionService().GetCallbackPayload(s.GetCtx(), msg.Arguments)
	if err != nil {
		logger.Warn("Invalid or expired callback token",
			"token", msg.Arguments, "err", err)
		_ = s.GetTgClient().SendMessage("Кнопка устарела, попробуйте снова", msg.UserID)
		return true, nil
	}

	// Проверка прав: тот ли пользователь нажал?
	if payload.UserID != msg.UserID {
		logger.Warn("User tried to access another user's task",
			"user_id", msg.UserID,
			"payload_user_id", payload.UserID,
			"task_id", payload.TaskID,
		)
		_ = s.GetTgClient().SendMessage("Доступ запрещён", msg.UserID)
		return true, nil
	}

	if payload.Action == "" {
		logger.Error("Unknown callback action", "action", payload.Action)
		return false, fmt.Errorf("unknown action: %s", payload.Action)
	}

	sessioData := make(map[string]any)

	// Сохраняем данные по задаче
	sessioData["payload_task_id"] = payload.TaskID
	if payload.Field != "" {
		sessioData["payload_field"] = payload.Field
	}

	_, _ = sessionutils.UpdateSessionData(s.GetCtx(), s.GetSessionService(), msg.UserID, sessioData)

	// Выполняем целевую команду
	if cmd, exists := s.GetCommandRegistry().GetCommand(payload.Action); exists {

		// TODO возможно session нужно реально будет записывать в сессию
		err = cmd.Execute(s, msg)
		if err != nil {
			return false, err
		}
		_, _ = sessionutils.ClearSessionData(s.GetCtx(), s.GetSessionService(), msg.UserID)

		return true, nil
	}

	return false, fmt.Errorf("command not found: %s", payload.Action)
}
