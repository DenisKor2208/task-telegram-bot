package commands

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnStart Команды стартовых действий.
var BtnStart = []bottypes.TgRowButtons{
	//{bottypes.TgInlineButton{DisplayName: "Добавить задачу", Value: "/start dfgdfgdf df df gdfgdfgdfgdf dfg 7 = 4"}},
	{bottypes.TgInlineButton{DisplayName: "Добавить задачу", Value: "/add_task"}},
	{bottypes.TgInlineButton{DisplayName: "Просмотреть задачи", Value: "/view_tasks"}},
	{bottypes.TgInlineButton{DisplayName: "Удалить задачу", Value: "/delete_task"}},
	{bottypes.TgInlineButton{DisplayName: "Редактировать задачу", Value: "/edit_task"}},
	{bottypes.TgInlineButton{DisplayName: "Выполнить задачу", Value: "/completed_task"}},
	{bottypes.TgInlineButton{DisplayName: "Завершить задачу", Value: "/closed_task"}},
}
