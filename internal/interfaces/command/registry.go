package command

// Registry — интерфейс для реестра команд.
type Registry interface {
	GetCommand(string) (Command, bool)
	RegisterCommand(string, Command)
}
