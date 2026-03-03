package views_tasks

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnFilterInProgress Команды для пункта меню "В процессе"
var BtnFilterInProgress = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{
		DisplayName: "На главную",
		Value:       "/start",
	}},
}
