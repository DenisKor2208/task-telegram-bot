// Package settimezone
package settimezone

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnSetTimezoneConfirm Команды при удалении задачи
var BtnSetTimezoneConfirm = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnMainMenu, Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnSettings, Value: "/settings_action"}},
}
