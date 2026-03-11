// Package storage
package storage

import (
	"context"

	"github.com/DenisKor2208/task-telegram-bot/internal/models"
)

type Status interface {
	GetStatusByID(context.Context, int) (*models.Status, error)
	//CreateStatus(context.Context, *models.Status) (bool, error)
}
