package commands

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
)

// RegistryCommands — простая реализация реестра.
type RegistryCommands struct {
	commands map[string]command.Command
}

// NewRegistryCommands — конструктор для реестра.
func NewRegistryCommands() *RegistryCommands {
	return &RegistryCommands{
		commands: make(map[string]command.Command),
	}
}

// RegisterCommand — метод для регистрации команды.
func (r *RegistryCommands) RegisterCommand(name string, cmd command.Command) {
	r.commands[name] = cmd
}

// GetCommand — метод для получения команды.
func (r *RegistryCommands) GetCommand(name string) (command.Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

// AllCommands возвращает все зарегистрированные команды.
// Нужно для поиска callback-обработчиков.
func (r *RegistryCommands) AllCommands() []command.Command {
	cmds := make([]command.Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}
