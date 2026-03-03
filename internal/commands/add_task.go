package commands

import (
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
)

const (
	COMMAND_ADD_TASKS = "/add_task"
)

// AddTaskCommand - структура команды /add_task - "Добавить задачу"
type AddTaskCommand struct{}

func (c *AddTaskCommand) Execute(s types.Model, msg types.Message, session *models.UserSession) error {
	// Определяем отображаемое имя: сначала UserDisplayName, иначе UserName.
	displayName := msg.UserDisplayName
	if len(displayName) == 0 {
		displayName = msg.UserName
	}

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTAddTask, displayName)
	textTemplate := resources.TXTAddTaskTemplate

	//_, _ = sessionutils.SetSessionStateFromCommand(s.GetCtx(), s.GetSessionService(), msg.UserID, COMMAND_ADD_TASKS)
	//if err != nil {
	//	return err
	//}

	// Отправляем сообщение без кнопок
	//return s.GetTgClient().SendMessage(text, msg.UserID)
	return s.GetTgClient().SendMessageWithTemplate(text, textTemplate, msg.UserID)
}
