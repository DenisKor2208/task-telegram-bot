// Package viewstasks
package viewstasks

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnBaseFilter Кнопки пункта меню при фильтрации задач по статусам
var BtnBaseFilter = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnMainMenu, Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnBack, Value: "/view_tasks"}},
}
