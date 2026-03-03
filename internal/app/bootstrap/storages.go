package bootstrap

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/jmoiron/sqlx"
)

type Storages struct {
	UserStorage   *repositories.UserStorage
	StatusStorage *repositories.StatusStorage
	TaskStorage   *repositories.TaskStorage
}

func StoragesFromDB(db *sqlx.DB) *Storages {
	return &Storages{
		UserStorage:   repositories.NewUserStorage(db),
		StatusStorage: repositories.NewStatusStorage(db),
		TaskStorage:   repositories.NewTaskStorage(db),
	}
}
