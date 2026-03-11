package commands

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/commands/set_timezone"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// SetTimezoneCommand - команда для выбора часового пояса "/set_timezone"
type SetTimezoneCommand struct{}

// timezoneOptions - предопределённый список часовых поясов с отображаемыми названиями.
var timezoneOptions = []struct {
	Display string
	Value   string
}{
	{"Москва (MSK)", "Europe/Moscow"},
	{"Санкт-Петербург (MSK)", "Europe/Moscow"},
	{"Калининград (UTC+2)", "Europe/Kaliningrad"},
	{"Самара (UTC+4)", "Europe/Samara"},
	{"Екатеринбург (UTC+5)", "Asia/Yekaterinburg"},
	{"Омск (UTC+6)", "Asia/Omsk"},
	{"Красноярск (UTC+7)", "Asia/Krasnoyarsk"},
	{"Иркутск (UTC+8)", "Asia/Irkutsk"},
	{"Якутск (UTC+9)", "Asia/Yakutsk"},
	{"Владивосток (UTC+10)", "Asia/Vladivostok"},
	{"Магадан (UTC+11)", "Asia/Magadan"},
	{"Камчатка (UTC+12)", "Asia/Kamchatka"},
}

func (c *SetTimezoneCommand) Execute(s command.Model, msg messaging.Message) error {
	var rows []bottypes.TgRowButtons
	for _, opt := range timezoneOptions {
		rows = append(rows, bottypes.TgRowButtons{
			{
				DisplayName: opt.Display,
				Value:       "/set_timezone_confirm " + opt.Value,
			},
		})
	}

	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := helpers.GetDisplayName(msg)

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTSetTimezone, displayName)

	return s.GetTgClient().ShowInlineButtons(text, rows, msg.UserID)
}

// NextStep Следующая команда
func (c *SetTimezoneCommand) NextStep() command.Command {
	return &set_timezone.SetTimezoneConfirmCommand{}
}

// InputField Поле для сохранения данных
func (c *SetTimezoneCommand) InputField() string {
	return ""
}
