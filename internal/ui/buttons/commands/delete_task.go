package commands

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnDeleteTask Команды для пункта меню "Удалить задачу"
var BtnDeleteTask = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnMainMenu, Value: "/start"}},
}
