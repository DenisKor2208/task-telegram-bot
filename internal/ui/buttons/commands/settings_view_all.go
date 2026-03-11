// Package commands
package commands

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnSettingsViewAll Команды стартовых действий.
var BtnSettingsViewAll = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: "На главную", Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: "Настройки", Value: "/settings_action"}},
}
