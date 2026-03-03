package commands

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnClosedTask Команды для пункта меню "Завершить задачу"
var BtnClosedTask = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{
		DisplayName: "На главную",
		Value:       "/start",
	}},
}
