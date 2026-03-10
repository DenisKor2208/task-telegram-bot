// Package bootstrap
package bootstrap

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/clients/tg"
	"github.com/DenisKor2208/task-telegram-bot/internal/config"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/services"
)

func DeadlineNotifier(
	cfg *config.Service,
	taskStorage *repositories.TaskStorage,
	userStorage *repositories.UserStorage,
	notifTypeStorage *repositories.NotificationTypeStorage,
	taskNotifStorage *repositories.TaskNotificationStorage,
	tgClient *tg.Client,
) *services.DeadlineNotifier {
	interval := cfg.GetConfig().NotificationCheckInterval
	return services.NewDeadlineNotifier(
		taskStorage,
		userStorage,
		notifTypeStorage,
		taskNotifStorage,
		tgClient,
		interval,
	)
}
