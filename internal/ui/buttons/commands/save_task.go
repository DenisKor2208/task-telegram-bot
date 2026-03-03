package commands

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnSaveTask Команды после сохранения задач
var BtnSaveTask = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: "На главную", Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: "Назад", Value: "/views_tasks"}},
}
