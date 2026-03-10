package services

import (
	"context"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
)

// OverdueChecker проверяет и обновляет статус просроченных задач.
type OverdueChecker struct {
	taskStorage *repositories.TaskStorage
	interval    time.Duration
}

// NewOverdueChecker создаёт новый экземпляр.
func NewOverdueChecker(taskStorage *repositories.TaskStorage, intervalSec int) *OverdueChecker {
	return &OverdueChecker{
		taskStorage: taskStorage,
		interval:    time.Duration(intervalSec) * time.Second,
	}
}

// Start запускает периодическую проверку в фоновой горутине.
// Завершается при отмене контекста.
func (oc *OverdueChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(oc.interval)
	defer ticker.Stop()

	// Сразу выполняем первую проверку
	oc.checkOverdueTasks(ctx)

	for {
		select {
		case <-ticker.C:
			oc.checkOverdueTasks(ctx)
		case <-ctx.Done():
			logger.Info("OverdueChecker остановлен")
			return
		}
	}
}

// checkOverdueTasks выполняет обновление статусов.
func (oc *OverdueChecker) checkOverdueTasks(ctx context.Context) {
	affected, err := oc.taskStorage.UpdateOverdueTasks(ctx)
	if err != nil {
		logger.Error("Ошибка обновления просроченных задач", "error", err)
		return
	}
	if affected > 0 {
		logger.Info("Обновлены просроченные задачи", "count", affected)
	}
}
