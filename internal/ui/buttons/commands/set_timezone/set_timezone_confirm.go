// Package set_timezone
package set_timezone

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnSetTimezoneConfirm Команды при удалении задачи
var BtnSetTimezoneConfirm = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: "На главную", Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: "Настройки", Value: "/settings_action"}},
}
