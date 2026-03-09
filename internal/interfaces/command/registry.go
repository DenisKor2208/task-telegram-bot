package command

// Registry — интерфейс для реестра команд.
type Registry interface {
	//GetCommand - Возвращает зарегистрированную команду по имени
	GetCommand(string) (Command, bool)

	//RegisterCommand - Регистрирует команду под определённым именем
	RegisterCommand(string, Command)

	//AllCommands - Возвращает все зарегистрированные команды
	AllCommands() []Command
}
