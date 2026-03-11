package commands

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// BtnOther
var BtnOther = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: resources.BtnMainMenu, Value: "/start"}},
}
