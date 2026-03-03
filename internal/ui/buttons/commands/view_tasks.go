package commands

import "github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"

// BtnViewTasks Команды для пункта меню "Просмотреть задачи"
var BtnViewTasks = []bottypes.TgRowButtons{
	{bottypes.TgInlineButton{DisplayName: "В процессе", Value: "/filter_in_progress"}},
	{bottypes.TgInlineButton{DisplayName: "Выполнено", Value: "/filter_completed"}},
	{bottypes.TgInlineButton{DisplayName: "Просрочено", Value: "/filter_overdue"}},
	{bottypes.TgInlineButton{DisplayName: "Завершено", Value: "/filter_closed"}},
	{bottypes.TgInlineButton{DisplayName: "Все задачи", Value: "/filter_all"}},
	{bottypes.TgInlineButton{DisplayName: "На главную", Value: "/start"}},
}
