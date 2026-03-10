package services

import (
	"context"
	"fmt"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
)

type DeadlineNotifier struct {
	taskStorage      *repositories.TaskStorage
	userStorage      *repositories.UserStorage
	notifTypeStorage *repositories.NotificationTypeStorage
	taskNotifStorage *repositories.TaskNotificationStorage
	tgClient         messaging.Sender
	interval         time.Duration
}

func NewDeadlineNotifier(
	taskStorage *repositories.TaskStorage,
	userStorage *repositories.UserStorage,
	notifTypeStorage *repositories.NotificationTypeStorage,
	taskNotifStorage *repositories.TaskNotificationStorage,
	tgClient messaging.Sender,
	intervalSec int,
) *DeadlineNotifier {
	return &DeadlineNotifier{
		taskStorage:      taskStorage,
		userStorage:      userStorage,
		notifTypeStorage: notifTypeStorage,
		taskNotifStorage: taskNotifStorage,
		tgClient:         tgClient,
		interval:         time.Duration(intervalSec) * time.Second,
	}
}

func (n *DeadlineNotifier) Start(ctx context.Context) {
	ticker := time.NewTicker(n.interval)
	defer ticker.Stop()

	n.checkDeadlines(ctx)

	for {
		select {
		case <-ticker.C:
			n.checkDeadlines(ctx)
		case <-ctx.Done():
			logger.Info("DeadlineNotifier остановлен")
			return
		}
	}
}

func (n *DeadlineNotifier) checkDeadlines(ctx context.Context) {
	// Получаем активные типы уведомлений
	notifTypes, err := n.notifTypeStorage.GetAllEnabled(ctx)
	if err != nil {
		logger.Error("DeadlineNotifier: не удалось получить типы уведомлений", "error", err)
		return
	}
	if len(notifTypes) == 0 {
		return
	}

	// Получаем задачи для уведомлений
	tasks, err := n.taskStorage.GetTasksForNotification(ctx)
	if err != nil {
		logger.Error("DeadlineNotifier: не удалось получить задачи", "error", err)
		return
	}

	now := time.Now()

	for _, task := range tasks {
		// Получаем уже отправленные типы для этой задачи
		sentTypes, err := n.taskNotifStorage.GetSentTypesForTask(ctx, task.ID)
		if err != nil {
			logger.Error("DeadlineNotifier: ошибка получения отправленных уведомлений", "task_id", task.ID, "error", err)
			continue
		}
		sentMap := make(map[int]bool, len(sentTypes))
		for _, id := range sentTypes {
			sentMap[id] = true
		}

		for _, nt := range notifTypes {
			if sentMap[nt.ID] {
				continue
			}

			// Время, когда нужно отправить уведомление
			notifyTime := task.Deadline.Add(-time.Duration(nt.TimeBefore) * time.Second)
			if !now.Before(notifyTime) {
				// Отправляем
				n.sendNotification(ctx, task, nt)
				// Записываем факт отправки (игнорируем дубликат)
				if err := n.taskNotifStorage.Insert(ctx, task.ID, nt.ID); err != nil {
					logger.Error("DeadlineNotifier: не удалось записать отправку", "task_id", task.ID, "type_id", nt.ID, "error", err)
				}
			}
		}
	}
}

func (n *DeadlineNotifier) sendNotification(ctx context.Context, task *models.Task, nt *models.NotificationType) {
	user, err := n.userStorage.GetUserByID(ctx, task.UserID)
	if err != nil {
		logger.Error("DeadlineNotifier: не удалось получить пользователя", "user_id", task.UserID, "error", err)
		return
	}

	deadlineStr := task.Deadline.Format("02.01.2006 15:04")
	text := fmt.Sprintf("🔔 Напоминание: %s\nДедлайн: %s\n⏰ %s", task.Description, deadlineStr, nt.Name)

	if err := n.tgClient.SendMessage(text, int64(user.TgID)); err != nil {
		logger.Error("DeadlineNotifier: ошибка отправки сообщения", "user_id", user.TgID, "error", err)
	} else {
		logger.Info("DeadlineNotifier: отправлено уведомление", "user_id", user.TgID, "task_id", task.ID, "type", nt.Name)
	}
}
