package models

import "time"

type TaskNotification struct {
	ID                 int       `db:"id"`
	TaskID             int       `db:"task_id"`              // ID задачи
	NotificationTypeID int       `db:"notification_type_id"` // ID типа уведомления
	SentAt             time.Time `db:"sent_at"`              // время отправки уведомления
}
