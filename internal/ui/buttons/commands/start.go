package commands

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnStart Кнопки пункта меню стартовых действий
var BtnStart = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnTasks, Value: "/tasks_action"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnSettings, Value: "/settings_action"}},
}
