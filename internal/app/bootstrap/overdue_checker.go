package bootstrap

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/config"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/services"
)

// OverdueChecker создаёт и возвращает экземпляр сервиса проверки просроченных задач.
func OverdueChecker(cfg *config.Service, taskStorage *repositories.TaskStorage) *services.OverdueChecker {
	interval := cfg.GetConfig().OverdueCheckInterval
	return services.NewOverdueChecker(taskStorage, interval)
}
