// Package commands
package commands

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnSettingsAction Команды стартовых действий.
var BtnSettingsAction = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: "Выбрать часовой пояс", Value: "/set_timezone"}},
	{bottypes.TgInlineButton{DisplayName: "Просмотреть текущие настройки", Value: "/settings_view_all"}},
	// {bottypes.TgInlineButton{DisplayName: "Просмотреть текущие настройки", Value: "/view_current_settings"}},
}
