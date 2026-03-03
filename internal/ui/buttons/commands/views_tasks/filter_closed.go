package views_tasks

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnFilterClosed Команды для пункта меню "Завершено"
var BtnFilterClosed = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{
		DisplayName: "На главную",
		Value:       "/start",
	}},
}
