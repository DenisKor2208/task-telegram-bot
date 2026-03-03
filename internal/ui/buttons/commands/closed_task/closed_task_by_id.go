package closed_task

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnClosedTaskByID Команды при удалении задачи
var BtnClosedTaskByID = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: "На главную", Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: "Назад", Value: "/closed_task"}},
}
