package bootstrap

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/config"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/dbutils"
	"github.com/jmoiron/sqlx"
)

func Database(cfg *config.Service) (*sqlx.DB, error) {
	db, err := dbutils.NewDBConnect(cfg.GetConfig().ConnectionStringDB)
	if err != nil {
		return nil, err
	}
	return db, nil
}
