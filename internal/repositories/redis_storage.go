package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/redis/go-redis/v9"
)

// RedisStorage предоставляет абстракцию для операций с Redis.
// Включает сериализацию JSON, таймауты и логирование.
type RedisStorage struct {
	client *redis.Client
}

// NewRedisStorage создаёт новый экземпляр RedisStorage.
// Принимает уже инициализированный клиент Redis.
func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

// Ping проверяет соединение с Redis.
// Возвращает результат пинга или ошибку.
func (rs *RedisStorage) Ping(ctx context.Context) (string, error) {
	// Используем переданный ctx, но добавляем таймаут, если его нет
	ctx, cancel := rs.applyDefaultTimeout(ctx)
	defer cancel()

	pong, err := rs.client.Ping(ctx).Result()
	if err != nil {
		logger.Error("Ошибка пинга Redis", "err", err)
		return "", fmt.Errorf("redis ping failed: %w", err)
	}
	logger.Info("Redis ping успешен", "pong", pong)
	return pong, nil
}

// Set сохраняет значение по ключу с опциональным TTL.
// value сериализуется в JSON. Если TTL == 0, ключ бессрочный.
func (rs *RedisStorage) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Сериализуем value в JSON
	data, err := json.Marshal(value)
	if err != nil {
		logger.Error("Ошибка сериализации значения для Redis", "key", key, "err", err)
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// Таймаут для операции
	ctx, cancel := rs.applyDefaultTimeout(ctx)
	defer cancel()

	err = rs.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		logger.Error("Ошибка сохранения в Redis", "key", key, "err", err)
		return fmt.Errorf("redis set failed: %w", err)
	}
	logger.Info("Значение сохранено в Redis", "key", key, "ttl", ttl)
	return nil
}

// Get получает значение по ключу и десериализует его из JSON.
// Возвращает ошибку, если ключ не найден (redis.Nil).
func (rs *RedisStorage) Get(ctx context.Context, key string, dest interface{}) error {
	// Таймаут для операции
	ctx, cancel := rs.applyDefaultTimeout(ctx)
	defer cancel()

	val, err := rs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		logger.Warn("Ключ не найден в Redis", "key", key)
		return redis.Nil
	}

	if err != nil {
		logger.Error("Ошибка получения из Redis", "key", key, "err", err)
		return fmt.Errorf("redis get failed: %w", err)
	}

	// Десериализуем из JSON
	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		logger.Error("Ошибка десериализации значения из Redis", "key", key, "err", err)
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}
	logger.Info("Значение получено из Redis", "key", key)
	return nil
}

// Del удаляет ключ(и) из Redis.
// Возвращает количество удалённых ключей.
func (rs *RedisStorage) Del(ctx context.Context, keys ...string) (int64, error) {
	// Таймаут для операции
	ctx, cancel := rs.applyDefaultTimeout(ctx)
	defer cancel()

	deleted, err := rs.client.Del(ctx, keys...).Result()
	if err != nil {
		logger.Error("Ошибка удаления из Redis", "keys", keys, "err", err)
		return 0, fmt.Errorf("redis del failed: %w", err)
	}
	logger.Info("Ключи удалены из Redis", "keys", keys, "deleted", deleted)
	return deleted, nil
}

// GetDel получает значение по ключу и удаляет его атомарно (Redis >= 6.2).
// Значение десериализуется в dest. Возвращает redis.Nil, если ключ не найден.
func (rs *RedisStorage) GetDel(ctx context.Context, key string, dest interface{}) error {
	ctx, cancel := rs.applyDefaultTimeout(ctx)
	defer cancel()

	// Используем настоящий GETDEL
	val, err := rs.client.GetDel(ctx, key).Result()
	if err == redis.Nil {
		logger.Warn("Ключ не найден в Redis (GetDel)", "key", key)
		return redis.Nil
	}
	if err != nil {
		logger.Error("Ошибка выполнения GetDel в Redis", "key", key, "err", err)
		return fmt.Errorf("redis getdel failed: %w", err)
	}

	// Десериализуем JSON
	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		logger.Error("Ошибка десериализации значения из Redis (GetDel)", "key", key, "err", err)
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	logger.Info("Значение получено и удалено в Redis", "key", key)
	return nil
}

// Exists проверяет существование ключа.
// Возвращает true, если ключ существует.
func (rs *RedisStorage) Exists(ctx context.Context, key string) (bool, error) {
	// Таймаут для операции
	ctx, cancel := rs.applyDefaultTimeout(ctx)
	defer cancel()

	count, err := rs.client.Exists(ctx, key).Result()
	if err != nil {
		logger.Error("Ошибка проверки существования ключа в Redis", "key", key, "err", err)
		return false, fmt.Errorf("redis exists failed: %w", err)
	}
	exists := count > 0
	logger.Info("Проверка существования ключа в Redis", "key", key, "exists", exists)
	return exists, nil
}

// Expire устанавливает TTL для существующего ключа.
// Возвращает true, если TTL установлен.
func (rs *RedisStorage) Expire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	// Таймаут для операции
	ctx, cancel := rs.applyDefaultTimeout(ctx)
	defer cancel()

	set, err := rs.client.Expire(ctx, key, ttl).Result()
	if err != nil {
		logger.Error("Ошибка установки TTL в Redis", "key", key, "ttl", ttl, "err", err)
		return false, fmt.Errorf("redis expire failed: %w", err)
	}
	logger.Info("TTL установлен для ключа в Redis", "key", key, "ttl", ttl, "set", set)
	return set, nil
}

// Close закрывает соединение с Redis.
// Вызывайте в defer в main.go.
func (rs *RedisStorage) Close() error {
	return rs.client.Close()
}

// applyDefaultTimeout применяет 5-секундный таймаут, если в контексте нет дедлайна.
func (rs *RedisStorage) applyDefaultTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); !ok {
		return context.WithTimeout(ctx, 5*time.Second)
	}
	// Дедлайн уже есть, ничего не делаем
	return ctx, func() {} // Возвращаем no-op cancel
}
