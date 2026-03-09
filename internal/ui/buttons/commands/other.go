package commands

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnOther
var BtnOther = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: "На главную", Value: "/start"}},
}
