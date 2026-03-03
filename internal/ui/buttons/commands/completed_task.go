package commands

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnCompletedTask Команды для пункта меню "Выполнить задачу"
var BtnCompletedTask = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{
		DisplayName: "На главную",
		Value:       "/start",
	}},
}
