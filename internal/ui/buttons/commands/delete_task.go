package commands

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnDeleteTask Команды для пункта меню "Удалить задачу"
var BtnDeleteTask = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{
		DisplayName: "На главную",
		Value:       "/start",
	}},
}
