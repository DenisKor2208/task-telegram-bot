package commands

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
)

// RegistryCommands — простая реализация реестра.
type RegistryCommands struct {
	commands map[string]types.Command
}

// NewRegistryCommands — конструктор для реестра.
func NewRegistryCommands() *RegistryCommands {
	return &RegistryCommands{
		commands: make(map[string]types.Command),
	}
}

// RegisterCommand — метод для регистрации команды.
func (r *RegistryCommands) RegisterCommand(name string, cmd types.Command) {
	r.commands[name] = cmd
}

// GetCommand — метод для получения команды.
func (r *RegistryCommands) GetCommand(name string) (types.Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}
