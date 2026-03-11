// Package deletetask
package deletetask

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnDeleteTaskByID Команды при удалении задачи
var BtnDeleteTaskByID = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnMainMenu, Value: "/start"}},
	{bottypes.TgInlineButton{DisplayName: resources.BtnBack, Value: "/delete_task"}},
}
