// Package bootstrap
package bootstrap

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/jmoiron/sqlx"
)

func TaskNotificationStorage(db *sqlx.DB) *repositories.TaskNotificationStorage {
	return repositories.NewTaskNotificationStorage(db)
}
