package bootstrap

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/redisutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/services"
)

func SessionService(redisClient *redisutils.RedisClient) *services.SessionService {
	redisStorage := repositories.NewRedisStorage(redisClient.Client)
	return services.NewSessionService(redisStorage, redisClient.SessionTTL())
}
