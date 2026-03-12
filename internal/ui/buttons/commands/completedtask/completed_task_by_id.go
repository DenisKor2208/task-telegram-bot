// Package completedtask
package completedtask

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnCompletedTaskByID Кнопки после смены статуса задачи на "Выполнено"
var BtnCompletedTaskByID = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnMainMenu, Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnBack, Value: "/completed_task"}},
}
