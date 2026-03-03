package bootstrap

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/clients/tg"
	"github.com/DenisKor2208/task-telegram-bot/internal/config"
)

func Telegram(cfg *config.Service) (*tg.Client, error) {
	handler := tg.HandlerFunc(tg.ProcessingMessages)
	client, err := tg.New(cfg, handler)
	if err != nil {
		return nil, err
	}
	return client, nil
}
