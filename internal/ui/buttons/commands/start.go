package commands

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnStart Команды стартовых действий.
var BtnStart = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: "Задачи", Value: "/tasks_action"}},
	{bottypes.TgInlineButton{DisplayName: "Настройки", Value: "/settings_action"}},
}
