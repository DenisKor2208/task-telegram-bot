package main

import (
	"context"

	"github.com/DenisKor2208/task-telegram-bot/internal/app/bootstrap"
	"github.com/DenisKor2208/task-telegram-bot/internal/clients/tg"
	"github.com/DenisKor2208/task-telegram-bot/internal/commands"
	"github.com/DenisKor2208/task-telegram-bot/internal/config"
	"github.com/DenisKor2208/task-telegram-bot/internal/helpers/redisutils"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/messages"
	"github.com/DenisKor2208/task-telegram-bot/internal/repositories"
	"github.com/DenisKor2208/task-telegram-bot/internal/services"

	"github.com/jmoiron/sqlx"
)

// App — центральная структура приложения.
// Содержит все зависимости и сервисы, необходимые для работы бота.
// Инициализируется через NewApp() и используется для запуска слушателя обновлений.
type App struct {
	cfg            *config.Service
	tgClient       *tg.Client
	dbConn         *sqlx.DB
	redisClient    *redisutils.RedisClient
	registry       *commands.RegistryCommands
	userStorage    *repositories.UserStorage
	statusStorage  *repositories.StatusStorage
	taskStorage    *repositories.TaskStorage
	sessionService *services.SessionService
	msgModel       *messages.Model
	storages       *bootstrap.Storages
}

// NewApp создаёт и инициализирует новое приложение.
//
// Принимает контекст для graceful shutdown и передачи в Model.
// Выполняет последовательную инициализацию всех компонентов.
// В случае ошибки на любом этапе возвращает nil и ошибку.
//
// Возвращает:
//   - *App: полностью инициализированное приложение.
//   - error: если один из этапов инициализации завершился с ошибкой.
func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	var err error

	// Загружает конфигурацию приложения
	if app.cfg, err = bootstrap.Config(); err != nil {
		logger.Fatal("Ошибка получения файла конфигурации:", "err", err)
		return nil, err
	}

	// Подключается к PostgreSQL с использованием строки подключения из конфига.
	if app.dbConn, err = bootstrap.Database(app.cfg); err != nil {
		logger.Fatal("Failed to connect to DB", "err", err)
		return nil, err
	}

	// Инициализирует клиент Redis с параметрами из конфигурации.
	// Используется для хранения сессий пользователей.
	if app.redisClient, err = bootstrap.Redis(app.cfg); err != nil {
		logger.Fatal("Failed to init Redis", "err", err)
		return nil, err
	}

	// Создаёт клиента Telegram API.
	if app.tgClient, err = bootstrap.Telegram(app.cfg); err != nil {
		logger.Fatal("Failed to init Telegram", "err", err)
		return nil, err
	}

	// Создаёт репозитории для работы с сущностями в БД:
	// - UserStorage
	// - StatusStorage
	// - TaskStorage
	// Все хранилища используют общее подключение a.dbConn.
	app.storages = bootstrap.StoragesFromDB(app.dbConn)
	app.userStorage = app.storages.UserStorage
	app.statusStorage = app.storages.StatusStorage
	app.taskStorage = app.storages.TaskStorage

	// Создаёт сервис управления сессиями.
	app.sessionService = bootstrap.SessionService(app.redisClient)

	app.registry = bootstrap.CommandRegistry()

	// создаёт основную модель обработки сообщений.
	app.msgModel = bootstrap.MessageModel(
		ctx,
		app.storages,
		app.tgClient,
		app.registry,
		app.sessionService,
	)

	return app, nil
}
