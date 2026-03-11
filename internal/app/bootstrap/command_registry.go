// Package bootstrap
package bootstrap

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/commands"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"

	cmdaddtask "github.com/DenisKor2208/task-telegram-bot/internal/commands/addtask"
	cmdclosedtask "github.com/DenisKor2208/task-telegram-bot/internal/commands/closed_task"
	cmdcompletedtask "github.com/DenisKor2208/task-telegram-bot/internal/commands/completed_task"
	cmddeletetask "github.com/DenisKor2208/task-telegram-bot/internal/commands/delete_task"
	cmdsettimezone "github.com/DenisKor2208/task-telegram-bot/internal/commands/set_timezone"
	cmdviewstasks "github.com/DenisKor2208/task-telegram-bot/internal/commands/viewstasks"
)

// CommandRegistry создаёт и наполняет реестр команд.
// Каждая команда (например, /add_task) связывается с конкретной реализацией Command.
func CommandRegistry() *commands.RegistryCommands {
	registry := commands.NewRegistryCommands()

	registry.RegisterCommand("start", &commands.StartCommand{})

	/***** Задачи ****/
	// Действия с задачами
	registry.RegisterCommand("tasks_action", &commands.TasksActionCommand{})

	// Добавить задачу
	registry.RegisterCommand("add_task", &commands.AddTaskCommand{})
	registry.RegisterCommand("save_task", &cmdaddtask.SaveTaskCommand{})

	// Просмотреть задачи
	registry.RegisterCommand("view_tasks", &commands.ViewTasksCommand{})
	registry.RegisterCommand("filter_all", &cmdviewstasks.BaseFilterCommand{
		Statuses:    nil, // все статусы
		ResourceKey: "filter_all",
	})
	registry.RegisterCommand("filter_in_progress", &cmdviewstasks.BaseFilterCommand{
		Statuses:    []int{repositories.StatusInProgress},
		ResourceKey: "filter_in_progress",
	})
	registry.RegisterCommand("filter_completed", &cmdviewstasks.BaseFilterCommand{
		Statuses:    []int{repositories.StatusCompleted},
		ResourceKey: "filter_completed",
	})
	registry.RegisterCommand("filter_overdue", &cmdviewstasks.BaseFilterCommand{
		Statuses:    []int{repositories.StatusOverdue},
		ResourceKey: "filter_overdue",
	})
	registry.RegisterCommand("filter_closed", &cmdviewstasks.BaseFilterCommand{
		Statuses:    []int{repositories.StatusClosed},
		ResourceKey: "filter_closed",
	})

	// Удалить задачу
	registry.RegisterCommand("delete_task", &commands.DeleteTaskCommand{})
	registry.RegisterCommand("delete_task_by_id", &cmddeletetask.DeleteTaskByIDCommand{})

	// Выполнить задачу
	registry.RegisterCommand("completed_task", &commands.CompletedTaskCommand{})
	registry.RegisterCommand("completed_task_by_id", &cmdcompletedtask.CompletedTaskByIDCommand{})

	// Завершить задачу
	registry.RegisterCommand("closed_task", &commands.ClosedTaskCommand{})
	registry.RegisterCommand("closed_task_by_id", &cmdclosedtask.ClosedTaskByIDCommand{})
	/********************/

	/***** Настройки ****/
	// Настройки пользователя
	registry.RegisterCommand("settings_action", &commands.SettingsActionCommand{})
	registry.RegisterCommand("settings_view_all", &commands.SettingsViewAllCommand{})

	// Часовой пояс пользователя
	registry.RegisterCommand("set_timezone", &commands.SetTimezoneCommand{})
	registry.RegisterCommand("set_timezone_confirm", &cmdsettimezone.SetTimezoneConfirmCommand{})
	/********************/

	return registry
}
