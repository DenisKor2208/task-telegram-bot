// Package settimezone
package settimezone

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnSetTimezoneConfirm Кнопки после сохранения выбранного часового пояса
var BtnSetTimezoneConfirm = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnMainMenu, Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnSettings, Value: "/settings_action"}},
}
