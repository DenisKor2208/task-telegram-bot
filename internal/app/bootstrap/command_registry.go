package bootstrap

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/commands"

	cmdaddtask "github.com/DenisKor2208/task-telegram-bot/internal/commands/add_task"
	cmdclosedtask "github.com/DenisKor2208/task-telegram-bot/internal/commands/closed_task"
	cmdcompletedtask "github.com/DenisKor2208/task-telegram-bot/internal/commands/completed_task"
	cmddeletetask "github.com/DenisKor2208/task-telegram-bot/internal/commands/delete_task"
	cmdedittasks "github.com/DenisKor2208/task-telegram-bot/internal/commands/edit_task"
	cmdviewstasks "github.com/DenisKor2208/task-telegram-bot/internal/commands/views_tasks"
)

// CommandRegistry создаёт и наполняет реестр команд.
// Каждая команда (например, /add_task) связывается с конкретной реализацией Command.
func CommandRegistry() *commands.RegistryCommands {
	registry := commands.NewRegistryCommands()

	registry.RegisterCommand("start", &commands.StartCommand{})
	registry.RegisterCommand("add_task", &commands.AddTaskCommand{})
	registry.RegisterCommand("save_task", &cmdaddtask.SaveTaskCommand{})

	// Просмотреть задачи
	registry.RegisterCommand("view_tasks", &commands.ViewTasksCommand{})
	registry.RegisterCommand("filter_all", &cmdviewstasks.FilterAllCommand{})
	registry.RegisterCommand("filter_in_progress", &cmdviewstasks.FilterInProgressCommand{})
	registry.RegisterCommand("filter_completed", &cmdviewstasks.FilterCompletedCommand{})
	registry.RegisterCommand("filter_overdue", &cmdviewstasks.FilterOverdueCommand{})
	registry.RegisterCommand("filter_closed", &cmdviewstasks.FilterClosedCommand{})

	// Удалить задачу
	registry.RegisterCommand("delete_task", &commands.DeleteTaskCommand{})
	registry.RegisterCommand("delete_task_by_id", &cmddeletetask.DeleteTaskByIdCommand{})

	// Редактировать задачу
	registry.RegisterCommand("edit_task", &commands.EditTaskCommand{})
	registry.RegisterCommand("edit_task_by_id", &cmdedittasks.EditTaskByIdCommand{})
	registry.RegisterCommand("update_task", &cmdedittasks.UpdateTaskCommand{})

	// Выполнить задачу
	registry.RegisterCommand("completed_task", &commands.CompletedTaskCommand{})
	registry.RegisterCommand("completed_task_by_id", &cmdcompletedtask.CompletedTaskByIdCommand{})

	// Завершить задачу
	registry.RegisterCommand("closed_task", &commands.ClosedTaskCommand{})
	registry.RegisterCommand("closed_task_by_id", &cmdclosedtask.ClosedTaskByIdCommand{})

	return registry
}
