package commands

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnViewTasks Кнопки пункта меню "Просмотреть задачи"
var BtnViewTasks = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnFilterInProgress, Value: "/filter_in_progress"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnFilterCompleted, Value: "/filter_completed"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnFilterOverdue, Value: "/filter_overdue"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnFilterClosed, Value: "/filter_closed"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnFilterAll, Value: "/filter_all"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnMainMenu, Value: "/start"}},
}
