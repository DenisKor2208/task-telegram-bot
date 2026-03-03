package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	configCommands "github.com/DenisKor2208/task-telegram-bot/internal/bot/config"
	"github.com/DenisKor2208/task-telegram-bot/internal/clients/tg"
	"github.com/DenisKor2208/task-telegram-bot/internal/commands"
	cmdaddtask "github.com/DenisKor2208/task-telegram-bot/internal/commands/add_task"
	cmdclosedtask "github.com/DenisKor2208/task-telegram-bot/internal/commands/closed_task"
	cmdcompletedtask "github.com/DenisKor2208/task-telegram-bot/internal/commands/completed_task"
	cmddeletetask "github.com/DenisKor2208/task-telegram-bot/internal/commands/delete_task"
	cmdedittasks "github.com/DenisKor2208/task-telegram-bot/internal/commands/edit_task"
	cmdviewstasks "github.com/DenisKor2208/task-telegram-bot/internal/commands/views_tasks"
	"github.com/DenisKor2208/task-telegram-bot/internal/config"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/dbutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/redisutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/messages"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/services"
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
)

func main() {
	logger.Info("Старт приложения")

	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("Ошибка получения файла конфигурации:", "err", err)
	}

	tgProcessingFuncHandler := tg.HandlerFunc(tg.ProcessingMessages)

	// Инициализация телеграм клиента.
	tgClient, err := tg.New(cfg, tgProcessingFuncHandler)
	if err != nil {
		logger.Fatal("Ошибка инициализации ТГ-клиента:", "err", err)
	}

	// Инициализация хранилищ (подключение к базе данных).
	dbconn, err := dbutils.NewDBConnect(cfg.GetConfig().ConnectionStringDB)
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных:", "err", err)
	}

	// Инициализация Redis клиента.
	redisClient, err := redisutils.NewRedisClient(cfg.GetConfig().Redis)
	if err != nil {
		logger.Fatal("Ошибка инициализации Redis:", "err", err)
	}
	defer redisClient.Close() // Graceful shutdown: закрываем клиент при выходе
	_ = redisClient           // Пока не используется — уберите, если передадите куда-то

	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGHUP,  //"Hang up" — обычно отправляется при закрытии терминала.
		syscall.SIGINT,  //"Interrupt" — отправляется при нажатии Ctrl+C в терминале.
		syscall.SIGTERM, //"Terminate" — стандартный сигнал для завершения процесса (от systemd, kill).
		syscall.SIGQUIT) //"Quit" — аналогично SIGINT, но может вызывать дамп стека (stack trace).
	defer cancel()

	registryCommands := commands.NewRegistryCommands()
	registryCommands.RegisterCommand("start", &commands.StartCommand{})
	registryCommands.RegisterCommand("add_task", &commands.AddTaskCommand{})
	registryCommands.RegisterCommand("save_task", &cmdaddtask.SaveTaskCommand{})

	// Просмотреть задачи
	registryCommands.RegisterCommand("view_tasks", &commands.ViewTasksCommand{})
	registryCommands.RegisterCommand("filter_all", &cmdviewstasks.FilterAllCommand{})
	registryCommands.RegisterCommand("filter_in_progress", &cmdviewstasks.FilterInProgressCommand{})
	registryCommands.RegisterCommand("filter_completed", &cmdviewstasks.FilterCompletedCommand{})
	registryCommands.RegisterCommand("filter_overdue", &cmdviewstasks.FilterOverdueCommand{})
	registryCommands.RegisterCommand("filter_closed", &cmdviewstasks.FilterClosedCommand{})

	// Удалить задачу
	registryCommands.RegisterCommand("delete_task", &commands.DeleteTaskCommand{})
	registryCommands.RegisterCommand("delete_task_by_id", &cmddeletetask.DeleteTaskByIdCommand{})

	registryCommands.RegisterCommand("edit_task", &commands.EditTaskCommand{})
	registryCommands.RegisterCommand("edit_task_by_id", &cmdedittasks.EditTaskByIdCommand{})
	registryCommands.RegisterCommand("update_task", &cmdedittasks.UpdateTaskCommand{})

	registryCommands.RegisterCommand("completed_task", &commands.CompletedTaskCommand{})
	registryCommands.RegisterCommand("completed_task_by_id", &cmdcompletedtask.CompletedTaskByIdCommand{})

	registryCommands.RegisterCommand("closed_task", &commands.ClosedTaskCommand{})
	registryCommands.RegisterCommand("closed_task_by_id", &cmdclosedtask.ClosedTaskByIdCommand{})

	userStorage := repositories.NewUserStorage(dbconn)
	statusStorage := repositories.NewStatusStorage(dbconn)
	taskStorage := repositories.NewTaskStorage(dbconn)
	redisStorage := repositories.NewRedisStorage(redisClient.Client)
	redisService := services.NewSessionService(redisStorage, redisClient.SessionTTL())

	// Инициализация конфигурации команд
	initalConfig := map[string]types.ConfigEntry{
		"start":              {"", []string{}},
		"view_tasks":         {"", []string{}},
		"filter_all":         {"", []string{}},
		"filter_closed":      {"", []string{}},
		"filter_completed":   {"", []string{}},
		"filter_in_progress": {"", []string{}},
		"filter_overdue":     {"", []string{}},
		"add_task": {
			"save_task",
			[]string{"title"},
		},
		"edit_task": {
			"edit_task_by_id",
			[]string{"title"},
		},
		"delete_task": {
			"",
			[]string{"title"},
		},
		"completed_task": {
			"",
			[]string{"title"},
		},
		"closed_task": {
			"",
			[]string{"title"},
		},
		"closed_task_by_id": {
			"",
			[]string{"task_id"},
		},
		"completed_task_by_id": {
			"",
			[]string{"task_id"},
		},
		"edit_task_by_id": {
			"update_task",
			[]string{"task_id", "title"},
		},
		"delete_task_by_id": {
			"",
			[]string{"task_id"},
		},
	}

	// Инициализация основной модели.
	msgModel := messages.New(ctx, userStorage, statusStorage, taskStorage, tgClient, registryCommands, redisService)
	configEntryService, err := configCommands.NewConfigEntryService(msgModel, initalConfig)
	if err != nil {
		logger.Fatal("Ошибка инициализации ConfigEntryService:", "err", err)
	}
	msgModel.SetConfigEntryService(configEntryService)

	log.Printf("Authorized on account %s", tgClient.Client.Self.UserName)

	// Запуск ТГ-клиента.
	tgClient.ListenUpdates(msgModel)

	// После выхода из ListenUpdates программа завершится (graceful shutdown).
	logger.Info("Application shutdown complete")
}
