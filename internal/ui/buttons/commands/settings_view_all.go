// Package commands
package commands

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnSettingsViewAll Кнопки пункта меню "Просмотреть текущие настройки"
var BtnSettingsViewAll = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnMainMenu, Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnSettings, Value: "/settings_action"}},
}
