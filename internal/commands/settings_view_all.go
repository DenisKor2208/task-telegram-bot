package commands

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	btnsettingsviewall "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands"
)

// SettingsViewAllCommand - команда для отображения текущих настроек пользователя "/settings_view_all"
type SettingsViewAllCommand struct{}

func (c *SettingsViewAllCommand) Execute(s command.Model, msg messaging.Message) error {
	// Получаем пользователя по Telegram ID
	user, err := s.GetUserStorage().GetUserByTgID(s.GetCtx(), int(msg.UserID))
	if err != nil {
		logger.Error("SettingsViewAllCommand: не удалось получить пользователя", "user_id", msg.UserID, "error", err)
		return s.GetTgClient().SendMessage("Не удалось загрузить настройки. Попробуйте позже.", msg.UserID)
	}

	// Формируем текст сообщения
	text := fmt.Sprintf("⚙️ <b>Ваши текущие настройки</b>\n\n🌍 Часовой пояс: <code>%s</code>", user.Timezone)

	// Инициализируем кнопки
	var buttons []bottypes.TgRowButtons

	buttons = append(buttons, btnsettingsviewall.BtnSettingsViewAll...)

	return s.GetTgClient().ShowInlineButtons(text, buttons, msg.UserID)
}

// NextStep Следующая команда
func (c *SettingsViewAllCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *SettingsViewAllCommand) InputField() string {
	return ""
}
