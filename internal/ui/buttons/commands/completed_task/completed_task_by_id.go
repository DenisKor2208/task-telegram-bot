package completed_task

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnCompletedTaskByID Команды при удалении задачи
var BtnCompletedTaskByID = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: "На главную", Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: "Назад", Value: "/completed_task"}},
}
