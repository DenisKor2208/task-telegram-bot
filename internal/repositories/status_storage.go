package repositories

import (
	"context"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/dbutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	StatusInProgress = 1 // Статус "В процессе"
	StatusCompleted  = 2 // Статус "Выполнено"
	StatusOverdue    = 3 // Статус "Просрочено"
	StatusClosed     = 4 // Статус "Завершено"
)

// StatusStorage represents storage for statuses
type StatusStorage struct {
	db *sqlx.DB
}

// NewStatusStorage returns new instance of StatusStorage
func NewStatusStorage(db *sqlx.DB) *StatusStorage {
	return &StatusStorage{db: db}
}

// GetStatusByID
func (ss *StatusStorage) GetStatusByID(ctx context.Context, statusId int) (*models.Status, error) {
	var status models.Status

	const sqlString = `SELECT * FROM statuses WHERE id = $1`

	// Выполнение запроса на получение данных.
	err := dbutils.Get(ctx, ss.db, &status, sqlString, statusId)
	if err != nil {
		return nil, err
	}

	return &status, nil
}
