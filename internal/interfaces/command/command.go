package command

import (
	"context"

	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/service"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/storage"
)

// Command — интерфейс для всех команд бота.
type Command interface {
	//Execute - Основная логика команды
	Execute(Model, messaging.Message) error
}

// Model — интерфейс для модели бота.
type Model interface {
	GetCtx() context.Context
	SetCtx(context.Context)
	GetCommandRegistry() Registry
	GetTgClient() messaging.Sender
	GetUserStorage() storage.User
	GetStatusStorage() storage.Status
	GetTaskStorage() storage.Task
	GetSessionService() service.Session
}
