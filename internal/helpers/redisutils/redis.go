package redisutils

import (
	"context"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/config"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/redis/go-redis/v9"
)

// RedisClient — обёртка над Redis, содержащая клиент и настройки.
type RedisClient struct {
	Client *redis.Client
	cfg    config.RedisConfig
}

// NewRedisClient создаёт и возвращает новый экземпляр RedisClient.
func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	// Парсим URL для корректной настройки клиента
	opt, err := redis.ParseURL(cfg.ConnectionString)
	if err != nil {
		logger.Error("Ошибка парсинга Redis URL: " + err.Error())
		return nil, err
	}

	client := redis.NewClient(opt)

	// Проверяем соединение с Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("Не удалось подключиться к Redis: " + err.Error())
		return nil, err
	}

	logger.Info("Подключено к Redis: " + opt.Addr)

	return &RedisClient{
		Client: client,
		cfg:    cfg,
	}, nil
}

// SessionTTL возвращает время жизни сессии из конфига (в секундах).
func (r *RedisClient) SessionTTL() time.Duration {
	return time.Duration(r.cfg.SessionTTL) * time.Second
}

// Close закрывает соединение с Redis.
func (r *RedisClient) Close() error {
	return r.Client.Close()
}
