package messages

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	configCommands "github.com/DenisKor2208/task-telegram-bot/internal/bot/config"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/sessionutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
	"github.com/pkg/errors"
)

// Model Модель бота (клиент, хранилище, последние команды пользователя)
type Model struct {
	ctx                context.Context
	userStorage        types.UserStorage
	statusStorage      types.StatusStorage
	taskStorage        types.TaskStorage
	tgClient           types.MessageSender   // Клиент
	commandRegistry    types.CommandRegistry // Хранилище команд
	lastUserCommand    map[int64]string      // Последняя выбранная пользователем команда
	sessionService     types.SessionService
	configEntryService types.ConfigEntryService
}

// New Генерация сущности для хранения клиента ТГ и хранилища пользователей.
func New(
	ctx context.Context,
	userStorage types.UserStorage,
	statusStorage types.StatusStorage,
	taskStorage types.TaskStorage,
	tgClient types.MessageSender,
	registry types.CommandRegistry,
	sessionService types.SessionService,
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

func (s *Model) GetTgClient() types.MessageSender {
	return s.tgClient
}

func (s *Model) GetCommandRegistry() types.CommandRegistry {
	return s.commandRegistry
}

func (s *Model) GetUserStorage() types.UserStorage {
	return s.userStorage
}

func (s *Model) GetStatusStorage() types.StatusStorage {
	return s.statusStorage
}

func (s *Model) GetTaskStorage() types.TaskStorage {
	return s.taskStorage
}

func (s *Model) GetLastUserCommand(userID int64) string {
	return s.lastUserCommand[userID]
}

func (s *Model) SetLastUserCommand(userID int64, command string) {
	s.lastUserCommand[userID] = command
}

func (s *Model) GetSessionService() types.SessionService {
	return s.sessionService
}

func (s *Model) GetConfigEntryService() types.ConfigEntryService {
	return s.configEntryService
}

func (s *Model) SetConfigEntryService(configEntryService *configCommands.ConfigEntryService) {
	s.configEntryService = configEntryService
}

// IncomingMessage Обработка входящего сообщения.
func (s *Model) IncomingMessage(msg types.Message) error {

	// Распознавание стандартных команд.
	if isNeedReturn, err := checkBotCommands(s, msg); err != nil || isNeedReturn {
		return err
	}

	// Отправка ответа по умолчанию.
	return s.tgClient.SendMessage(resources.TXTUnknownCommand, msg.UserID)
}

// Распознавание стандартных команд бота.
func checkBotCommands(s types.Model, msg types.Message) (bool, error) {

	// Если команда или аргументы не указаны, то не обрабатываем.
	if msg.Command == "" && msg.Arguments == "" {
		return false, errors.New("message is not a command or arguments")
	}

	// Получаем сессию пользователя
	session, err := s.GetSessionService().GetSession(s.GetCtx(), msg.UserID)
	if err != nil {
		logger.Error("Error getting user session", "err", err, "userID", msg.UserID)
		return false, err
	}
	// Если не команда, без аргумента и в сессии нету токена задачи
	if !msg.IsCommand && msg.Arguments != "" && session.State != "" {
		msg.Arguments = session.State
	}

	// Проверяем, не является ли аргумент токеном задачи
	if s.GetSessionService().IsCallbackToken(msg.Arguments) {
		return handleCallbackToken(s, msg, session)
	}

	// Обработка команд
	if msg.IsCommand {
		return handleCommand(s, msg, session)
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
func handleCommand(s types.Model, msg types.Message, session *models.UserSession) (bool, error) {

	var err error

	for cmdName, config := range s.GetConfigEntryService().GetConfig() {
		if msg.Command == cmdName {
			if cmd, exists := s.GetCommandRegistry().GetCommand(cmdName); exists {

				nextState := config.TargetCommand

				// Устанавливаем в сессию потенциальную следующую команду
				if nextState != "" {
					_, _ = sessionutils.UpdateSessionData(s.GetCtx(), s.GetSessionService(), msg, map[string]interface{}{"next_command": nextState})
				}

				err = cmd.Execute(s, msg, session)
				if err != nil {
					return false, err
				}

				// Сохраняем текущую выполненную команду
				_, _ = sessionutils.UpdateSessionData(s.GetCtx(), s.GetSessionService(), msg, map[string]interface{}{"last_command": cmdName})

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
func handleFSM(s types.Model, msg types.Message, session *models.UserSession) (bool, error) {

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
	session, _ = sessionutils.UpdateSessionData(
		s.GetCtx(), s.GetSessionService(), msg,
		map[string]interface{}{config.CommandFields[0]: msg.Arguments},
	)

	// Выполняем целевую команду
	err = cmd.Execute(s, msg, session)
	if err != nil {
		return false, err
	}

	_, _ = sessionutils.SetSessionTaskToken(s.GetCtx(), s.GetSessionService(), msg.UserID, "")
	_, _ = sessionutils.ClearSessionData(s.GetCtx(), s.GetSessionService(), msg.UserID)
	_, _ = sessionutils.UpdateSessionData(s.GetCtx(), s.GetSessionService(), msg, map[string]interface{}{"last_command": nextCommand})

	// Некоторых команд удаляем сессию полностью (диалог окончен)
	/* 	if slices.Contains([]string{"save_task", "update_task"}, config.TargetCommand) {
		_ = s.GetSessionService().DeleteSession(s.GetCtx(), msg.UserID)
	} */

	return true, nil
}

// handleCallbackToken — обрабатывает нажатие кнопки с callback-токеном
func handleCallbackToken(s types.Model, msg types.Message, session *models.UserSession) (bool, error) {
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

	// Сохраняем данные по задаче
	session.Data["payload_task_id"] = payload.TaskID
	if payload.Field != "" {
		session.Data["payload_field"] = payload.Field
	}

	// Выполняем целевую команду
	if cmd, exists := s.GetCommandRegistry().GetCommand(payload.Action); exists {

		// TODO возможно session нужно реально будет записывать в сессию
		err = cmd.Execute(s, msg, session)
		if err != nil {
			return false, err
		}
		_, _ = sessionutils.ClearSessionData(s.GetCtx(), s.GetSessionService(), msg.UserID)

		return true, nil
	}

	return false, fmt.Errorf("command not found: %s", payload.Action)
}
