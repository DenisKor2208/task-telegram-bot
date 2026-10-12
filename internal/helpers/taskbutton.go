package helpers

import (
	"fmt"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
)

// taskStatusView - как задача с данным статусом выглядит в списке:
// эмодзи в начале текста (для старых клиентов Telegram, где цвет не виден) и цвет кнопки.
type taskStatusView struct {
	emoji string
	style string
}

var taskStatusViews = map[int]taskStatusView{
	repositories.StatusInProgress: {emoji: "🔵", style: bottypes.ButtonStylePrimary},
	repositories.StatusCompleted:  {emoji: "🟢", style: bottypes.ButtonStyleSuccess},
	repositories.StatusOverdue:    {emoji: "🔴", style: bottypes.ButtonStyleDanger},
	repositories.StatusClosed:     {emoji: "⚪", style: ""}, // обычная кнопка
}

// TaskButton возвращает кнопку задачи для списков: эмодзи и цвет по статусу, описание, дедлайн
// в часовом поясе пользователя. По нажатию отправляется команда "<command> <id задачи>".
func TaskButton(task *models.Task, loc *time.Location, command string) bottypes.TgInlineButton {
	view := taskStatusViews[task.StatusID] // неизвестный статус — без эмодзи и цвета

	text := fmt.Sprintf("%s (%s)", task.Description, DeadlineLabel(task.Deadline, loc))
	if view.emoji != "" {
		text = view.emoji + " " + text
	}

	return bottypes.TgInlineButton{
		DisplayName: text,
		Value:       fmt.Sprintf("%s %d", command, task.ID),
		Style:       view.style,
	}
}
