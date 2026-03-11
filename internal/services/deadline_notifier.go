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
	if len(tasks) == 0 {
		return
	}

	// Собираем ID задач для массового запроса
	taskIDs := make([]int, 0, len(tasks))
	for _, task := range tasks {
		taskIDs = append(taskIDs, task.ID)
	}

	// Получаем мапу отправленных уведомлений для всех задач одним запросом
	sentMap, err := n.taskNotifStorage.GetSentMapForTasks(ctx, taskIDs)
	if err != nil {
		logger.Error("DeadlineNotifier: не удалось получить отправленные уведомления", "error", err)
		// Продолжаем, но для задач без данных sentMap будет пусто
	}

	// Рекомендуется использовать UTC для единообразия с хранимыми временами
	now := time.Now().UTC()

	for _, task := range tasks {
		// Получаем уже отправленные типы для этой задачи
		sentForTask := sentMap[task.ID] // если ключа нет, вернётся nil
		sent := make(map[int]bool, len(sentForTask))
		for _, id := range sentForTask {
			sent[id] = true
		}

		for _, nt := range notifTypes {
			if sent[nt.ID] {
				continue
			}

			// Время, когда нужно отправить уведомление
			notifyTime := task.Deadline.Add(-time.Duration(nt.TimeBefore) * time.Second)
			if !now.Before(notifyTime) {
				// Отправляем
				n.sendNotification(ctx, task, nt)
				// Записываем факт отправки (дубликаты игнорируются на уровне БД)
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

	// Загружаем локацию пользователя
	loc, err := time.LoadLocation(user.Timezone)
	if err != nil {
		logger.Warn("DeadlineNotifier: не удалось загрузить часовой пояс, используется UTC", "user_id", user.ID, "timezone", user.Timezone)
		loc = time.UTC
	}

	// Преобразуем время дедлайна в локальное время пользователя
	deadlineLocal := task.Deadline.In(loc)
	deadlineStr := deadlineLocal.Format("02.01.2006 15:04")

	text := fmt.Sprintf("🔔 Напоминание: %s\nДедлайн: %s (%s)\n⏰ %s", task.Description, deadlineStr, user.Timezone, nt.Name)

	if err := n.tgClient.SendMessage(text, int64(user.TgID)); err != nil {
		logger.Error("DeadlineNotifier: ошибка отправки сообщения", "user_id", user.TgID, "error", err)
	} else {
		logger.Info("DeadlineNotifier: отправлено уведомление", "user_id", user.TgID, "task_id", task.ID, "type", nt.Name)
	}
}
