// Package commands
package commands

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnTasksAction Кнопки пункта меню "Задачи"
var BtnTasksAction = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnAddTask, Value: "/add_task"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnViewTasks, Value: "/view_tasks"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnDeleteTask, Value: "/delete_task"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnCompleteTask, Value: "/completed_task"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnCloseTask, Value: "/closed_task"}},
}
