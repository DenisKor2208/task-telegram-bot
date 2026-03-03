package edit_task

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnUpdateTask Команды после сохранения задач
var BtnUpdateTask = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: "На главную", Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: "Назад", Value: "/edit_task"}},
}
