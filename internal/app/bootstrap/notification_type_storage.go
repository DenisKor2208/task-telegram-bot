// Package bootstrap
package bootstrap

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/jmoiron/sqlx"
)

func NotificationTypeStorage(db *sqlx.DB) *repositories.NotificationTypeStorage {
	return repositories.NewNotificationTypeStorage(db)
}
