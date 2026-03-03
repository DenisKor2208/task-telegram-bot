package config

import (
	"fmt"
	"sync"

	/* 	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/sessionutils"
	   	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	   	"github.com/DenisKor2208/task-telegram-bot/internal/models" */
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
	"github.com/pkg/errors"
)

// ConfigEntryService — сервис для работы с types.ConfigEntry и обработки сессий/команд.
type ConfigEntryService struct {
	model  types.Model
	config map[string]types.ConfigEntry
	mu     sync.RWMutex
}

// NewConfigEntryService — конструктор сервиса.
// Принимает зависимость model и начальную конфигурацию (может быть пустой map).
// Валидирует конфиг и возвращает ошибку, если она некорректна.
func NewConfigEntryService(model types.Model, initialConfig map[string]types.ConfigEntry) (*ConfigEntryService, error) {
	if model == nil {
		return nil, errors.New("model interface cannot be nil")
	}
	if initialConfig == nil {
		initialConfig = make(map[string]types.ConfigEntry)
	}
	/* 	if err := validateConfig(initialConfig); err != nil {
		return nil, err
	} */
	return &ConfigEntryService{
		model:  model,
		config: initialConfig,
	}, nil
}

func (ces *ConfigEntryService) GetConfig() map[string]types.ConfigEntry {
	ces.mu.RLock()
	defer ces.mu.RUnlock()
	return ces.config
}

// validateConfig проверяет корректность конфига.
// Вызывается в конструкторе.
func validateConfig(config map[string]types.ConfigEntry) error {
	for key, entry := range config {
		if entry.TargetCommand != "" && len(entry.CommandFields) > 0 {
			return fmt.Errorf("config entry for key '%s' has both TargetCommand and CommandFields, which is invalid", key)
		}
		if entry.TargetCommand != "" && entry.CommandFields == nil {
			return fmt.Errorf("config entry for key '%s' has TargetCommand but missing CommandFields ", key)
		}
		if len(entry.CommandFields) > 0 && (entry.TargetCommand != "" || entry.CommandFields != nil) {
			return fmt.Errorf("config entry for key '%s' has CommandFields but also  fields, which is invalid", key)
		}
	}
	return nil
}
