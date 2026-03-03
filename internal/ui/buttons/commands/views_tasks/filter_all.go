package views_tasks

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnFilterAll Команды для пункта меню "Все задачи"
var BtnFilterAll = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{
		DisplayName: "На главную",
		Value:       "/start",
	}},
}
