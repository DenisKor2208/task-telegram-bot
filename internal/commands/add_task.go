package commands

import (
	"context"
	"fmt"

	"github.com/DenisKor2208/task-telegram-bot/internal/commands/addtask"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
)

// AddTaskCommand - структура команды /add_task - "Добавить задачу"
type AddTaskCommand struct{}

// Execute — реализация команды /add_task.
func (c *AddTaskCommand) Execute(s command.Model, msg messaging.Message) error {
	return c.Prompt(s.GetCtx(), s, msg)
}

func (c *AddTaskCommand) Prompt(ctx context.Context, model command.Model, msg messaging.Message) error {
	displayName := helpers.GetDisplayName(msg)

	// Формируем текст приветствия
	text := fmt.Sprintf(resources.TXTAddTask, displayName)
	textTemplate := resources.TXTAddTaskTemplate

	return model.GetTgClient().SendMessageWithTemplate(text, textTemplate, msg.UserID)
}

// NextStep Следующая команда
func (c *AddTaskCommand) NextStep() command.Command {
	return &addtask.SaveTaskCommand{}
}

// InputField Поле для сохранения данных
func (c *AddTaskCommand) InputField() string {
	return "title"
}
