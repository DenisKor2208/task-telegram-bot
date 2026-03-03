package views_tasks

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnFilterCompleted Команды для пункта меню "Выполнено"
var BtnFilterCompleted = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{
		DisplayName: "На главную",
		Value:       "/start",
	}},
}
