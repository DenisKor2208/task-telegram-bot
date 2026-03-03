package types

import (
	"context"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/callbacktokenpayloadutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
)

// Command — интерфейс для всех команд бота.
type Command interface {
	Execute(Model, Message, *models.UserSession) error
}

// CommandRegistry — интерфейс для реестра команд.
type CommandRegistry interface {
	GetCommand(string) (Command, bool)
	RegisterCommand(string, Command)
}

type UserStorage interface {
	GetUserByTgID(context.Context, int) (*models.User, error)
	CreateUser(context.Context, *models.User) (*models.User, error)
}

type StatusStorage interface {
	GetStatusByID(context.Context, int) (*models.Status, error)
	//CreateStatus(context.Context, *models.Status) (bool, error)
}

type TaskStorage interface {
	GetTaskByID(context.Context, int64) (*models.Task, error)
	CreateTask(context.Context, *models.Task) (bool, error)
	GetAllTasks(context.Context) ([]*models.Task, error)
	GetTasksByStatusID(context.Context, []int) ([]*models.Task, error)
	DeleteTaskByID(context.Context, int64) error
	UpdateTask(context.Context, *models.Task) error
}

// Model — интерфейс для модели бота.
type Model interface {
	GetCtx() context.Context
	SetCtx(context.Context)
	GetTgClient() MessageSender
	GetCommandRegistry() CommandRegistry
	GetUserStorage() UserStorage
	GetStatusStorage() StatusStorage
	GetTaskStorage() TaskStorage
	GetLastUserCommand(int64) string
	SetLastUserCommand(int64, string)
	GetSessionService() SessionService
	GetConfigEntryService() ConfigEntryService
}

// Message Структура сообщения для обработки.
type Message struct {
	Text            string
	Command         string
	Arguments       string
	IsCommand       bool
	UserID          int64
	UserName        string
	UserDisplayName string
	Date            int64
	IsCallback      bool
	CallbackMsgID   string
}

// MessageSender Интерфейс для работы с сообщениями.
type MessageSender interface {
	SendMessage(text string, userID int64) error
	ShowInlineButtons(text string, buttons []bottypes.TgRowButtons, userID int64) error
	SendMessageWithTemplate(text string, textTemplate string, userID int64) error
}

type SessionService interface {
	GetSession(context.Context, int64) (*models.UserSession, error)
	SaveSession(context.Context, *models.UserSession) error
	DeleteSession(context.Context, int64) error

	// Методы для callback-токенов
	CreateCallbackToken(ctx context.Context, payload callbacktokenpayloadutils.CallbackTokenPayload, ttl time.Duration) (string, error)
	GetCallbackPayload(ctx context.Context, token string) (*callbacktokenpayloadutils.CallbackTokenPayload, error)
	IsCallbackToken(token string) bool
}

type ConfigEntryService interface {
	GetConfig() map[string]ConfigEntry
}

type ConfigEntry struct {
	TargetCommand string
	CommandFields []string
}
