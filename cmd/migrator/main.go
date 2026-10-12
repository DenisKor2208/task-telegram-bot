package main

// Для работы с миграциями

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/config"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {

	if len(os.Args) < 2 {
		helpInfo()
		os.Exit(2)
	}

	command := os.Args[1]

	if command == "help" || command == "--help" || command == "-h" {
		helpInfo()
		os.Exit(0)
	}

	// Создание миграции
	if command == "create" {
		if len(os.Args) < 3 {
			usageError("Укажите имя миграции: create <name>")
		}
		createMigration(os.Args[2])
		return
	}

	migrationsPath := "file://./migrations"

	logger.Info("Старт приложения")

	cfg, err := config.NewDatabase()
	if err != nil {
		logger.Fatal("Ошибка чтения настроек подключения к БД", "err", err)
	}

	// Подключение к БД
	db, err := sql.Open("postgres", cfg.ConnectionString)
	if err != nil {
		logger.Fatal("Ошибка подключения к БД", "err", err)
	}
	defer db.Close()

	// Настройка драйвера миграций
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Fatal("Ошибка настройки драйвера", "err", err)
	}

	// Источник миграций из файлов
	source, err := (&file.File{}).Open(migrationsPath)
	if err != nil {
		logger.Fatal("Ошибка открытия источника миграций", "err", err)
	}

	// Создание экземпляра migrate
	m, err := migrate.NewWithInstance("file", source, "postgres", driver)
	if err != nil {
		logger.Fatal("Ошибка создания migrate", "err", err)
	}
	defer m.Close()

	// Обработка команд
	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			logger.Fatal("Ошибка применения миграций", "err", err)
		}
		fmt.Println("Миграции применены!")
	case "down":
		if len(os.Args) < 3 {
			usageError("Укажите количество миграций для отката: down <n>")
		}
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			usageError("Количество миграций должно быть числом: down <n>")
		}
		if err := m.Steps(-n); err != nil && err != migrate.ErrNoChange {
			logger.Fatal("Ошибка отката миграций", "err", err)
		}
		fmt.Println("Миграции откатаны!")
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			logger.Fatal("Ошибка получения версии", "err", err)
		}
		fmt.Printf("Версия: %d, Dirty: %t\n", version, dirty)
	case "goto":
		if len(os.Args) < 3 {
			usageError("Укажите версию: goto <version>")
		}
		version, err := strconv.ParseUint(os.Args[2], 10, 64)
		if err != nil {
			usageError("Версия должна быть числом: goto <version>")
		}
		if err := m.Migrate(uint(version)); err != nil && err != migrate.ErrNoChange {
			logger.Fatal("Ошибка перехода к версии", "err", err)
		}
		fmt.Println("Переход выполнен!")
	default:
		usageError("Неизвестная команда: " + command)
	}
}

// usageError выводит ошибку использования (неверные аргументы) и справку, завершает работу с кодом 2.
// Без лога и стека вызовов: это подсказка пользователю, а не сбой программы.
func usageError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Fprintln(os.Stderr)
	helpInfo()
	os.Exit(2)
}

// helpInfo выводит справочную информацию
func helpInfo() {
	fmt.Println("Использование: migrate <command> [args...]")
	fmt.Println("Команды:")
	fmt.Println("  create <name>    - Создать новую миграцию")
	fmt.Println("  up               - Применить все миграции")
	fmt.Println("  down <n>         - Откатить n миграций")
	fmt.Println("  version          - Показать текущую версию")
	fmt.Println("  goto <version>   - Перейти к конкретной версии")
}

// createMigration Функция для создания новой миграции (up и down файлов)
func createMigration(name string) {
	timestamp := time.Now().Format("20060102150405") // Формат: YYYYMMDDHHMMSS
	upFile := fmt.Sprintf("migrations/%s_%s.up.sql", timestamp, name)
	downFile := fmt.Sprintf("migrations/%s_%s.down.sql", timestamp, name)

	// Создание up-файла
	if err := os.WriteFile(upFile, []byte("-- Миграция вверх: "+name+"\n\n"), 0644); err != nil {
		logger.Fatal("Ошибка создания up-файла", "err", err)
	}

	// Создание down-файла
	if err := os.WriteFile(downFile, []byte("-- Миграция вниз: "+name+"\n\n"), 0644); err != nil {
		logger.Fatal("Ошибка создания down-файла", "err", err)
	}

	fmt.Printf("Миграция создана: %s\n", upFile)
	fmt.Printf("Миграция создана: %s\n", downFile)
}
