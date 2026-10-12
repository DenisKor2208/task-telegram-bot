package commands

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnViewTasks Кнопки пункта меню "Просмотреть задачи".
// Цвета фильтров совпадают с цветами кнопок задач в этих статусах (см. helpers.TaskButton).
var BtnViewTasks = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnFilterInProgress, Value: "/filter_in_progress", Style: bottypes.ButtonStylePrimary}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnFilterCompleted, Value: "/filter_completed", Style: bottypes.ButtonStyleSuccess}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnFilterOverdue, Value: "/filter_overdue", Style: bottypes.ButtonStyleDanger}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnFilterClosed, Value: "/filter_closed"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnFilterAll, Value: "/filter_all"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnMainMenu, Value: "/start"}},
}
