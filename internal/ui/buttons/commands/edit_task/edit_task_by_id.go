package edit_task

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnEditTaskByID Команды при редактировании задачи
var BtnEditTaskByID = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{
		DisplayName: "На главную",
		Value:       "/start",
	}},
}
