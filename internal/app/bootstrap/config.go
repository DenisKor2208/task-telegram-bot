package bootstrap

import "github.com/DenisKor2208/task-telegram-bot/internal/config"

func Config() (*config.Service, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
