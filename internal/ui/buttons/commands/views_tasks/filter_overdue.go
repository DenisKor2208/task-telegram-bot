package views_tasks

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnFilterOverdue Команды для пункта меню "Просрочено"
var BtnFilterOverdue = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{
		DisplayName: "На главную",
		Value:       "/start",
	}},
}
