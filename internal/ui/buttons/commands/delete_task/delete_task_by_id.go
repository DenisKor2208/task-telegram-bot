// Package delete_task
package delete_task

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnDeleteTaskByID Команды при удалении задачи
var BtnDeleteTaskByID = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: "На главную", Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: "Назад", Value: "/delete_task"}},
}
