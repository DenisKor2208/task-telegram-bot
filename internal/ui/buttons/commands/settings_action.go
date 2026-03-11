// Package commands
package commands

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnSettingsAction Команды стартовых действий.
var BtnSettingsAction = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnSetTimezone, Value: "/set_timezone"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnViewSettings, Value: "/settings_view_all"}},
	// {bottypes.TgInlineButton{DisplayName: resources.BtnViewSettings, Value: "/view_current_settings"}},
}
