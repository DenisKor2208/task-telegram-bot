package config

// ConfigEntry — универсальная конфигурация для FSM-переходов и команд.
// - Для FSM: заполняем FSMTargetCommand и FSMFieldToUpdate.
// - Для команд: заполняем CommandFields (список полей для обновления).
// - Пустые поля игнорируются в логике.
type ConfigEntry struct {
	FSMTargetCommand string   // Команда, которая запустится, если есть только аргументы и предыдущее состояние совпадает (для FSM)
	FSMFieldToUpdate string   // Поле в сессии (UserSession.Data) для записи полного args (для FSM)
	CommandFields    []string // Список полей в сессии (UserSession.Data) для записи данных из args (для команд)
}
