// Package closedtask
package closedtask

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnClosedTaskByID Кнопки после смены статуса задачи на "Завершено"
var BtnClosedTaskByID = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnMainMenu, Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnBack, Value: "/closed_task"}},
}
