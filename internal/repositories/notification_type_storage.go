package repositories

import (
	"context"
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/dbutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/jmoiron/sqlx"
)

type NotificationTypeStorage struct {
	db *sqlx.DB
}

func NewNotificationTypeStorage(db *sqlx.DB) *NotificationTypeStorage {
	return &NotificationTypeStorage{db: db}
}

// GetAllEnabled возвращает все активные типы уведомлений.
func (s *NotificationTypeStorage) GetAllEnabled(ctx context.Context) ([]*models.NotificationType, error) {
	var types []*models.NotificationType
	query := `SELECT id, name, time_before, enabled, created_at, updated_at 
              FROM notification_types WHERE enabled = true ORDER BY time_before DESC`
	err := dbutils.Select(ctx, s.db, &types, query)
	if err != nil {
		return nil, fmt.Errorf("GetAllEnabled: %w", err)
	}
	return types, nil
}
