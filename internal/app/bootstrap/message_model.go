package bootstrap

import (
	"context"
	"github.com/DenisKor2208/task-telegram-bot/internal/clients/tg"
	"github.com/DenisKor2208/task-telegram-bot/internal/commands"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/messages"
	"github.com/DenisKor2208/task-telegram-bot/internal/services"
)

func MessageModel(
	ctx context.Context,
	storages *Storages,
	tgClient *tg.Client,
	registry *commands.RegistryCommands,
	sessionService *services.SessionService,
) *messages.Model {
	return messages.New(
		ctx,
		storages.UserStorage,
		storages.StatusStorage,
		storages.TaskStorage,
		tgClient,
		registry,
		sessionService,
	)
}
