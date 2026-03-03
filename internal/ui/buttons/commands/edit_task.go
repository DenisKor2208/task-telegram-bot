package commands

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnEditTask Команды для пункта меню "Редактировать задачу"
var BtnEditTask = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{
		DisplayName: "На главную",
		Value:       "/start",
	}},
}
