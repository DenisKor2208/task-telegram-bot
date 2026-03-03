package bootstrap

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/config"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/redisutils"
)

func Redis(cfg *config.Service) (*redisutils.RedisClient, error) {
	client, err := redisutils.NewRedisClient(cfg.GetConfig().Redis)
	if err != nil {
		return nil, err
	}
	return client, nil
}
