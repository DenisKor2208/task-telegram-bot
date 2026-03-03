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
	_ "github.com/lib/pq"
)

func main() {

	if len(os.Args) < 2 {
		helpInfo()
	}

	command := os.Args[1]

	if command == "help" || command == "--help" || command == "-h" {
		helpInfo()
		os.Exit(1)
	}

	migrationsPath := "file://./migrations"

	logger.Info("Старт приложения")

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("Ошибка получения файла конфигурации:", "err", err)
	}

	// Подключение к БД
	db, err := sql.Open("postgres", cfg.GetConfig().ConnectionStringDB)
	if err != nil {
		logger.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// Настройка драйвера миграций
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Fatal("Ошибка настройки драйвера:", err)
	}

	// Источник миграций из файлов
	source, err := (&file.File{}).Open(migrationsPath)
	if err != nil {
		logger.Fatal("Ошибка открытия источника миграций:", err)
	}

	// Создание экземпляра migrate
	m, err := migrate.NewWithInstance("file", source, "postgres", driver)
	if err != nil {
		logger.Fatal("Ошибка создания migrate:", err)
	}
	defer m.Close()

	// Обработка команд
	switch command {
	case "create":
		if len(os.Args) < 3 {
			logger.Fatal("Укажите имя миграции: create <name>")
		}
		name := os.Args[2]
		createMigration(name)
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			logger.Fatal("Ошибка применения миграций:", err)
		}
		fmt.Println("Миграции применены!")
	case "down":
		if len(os.Args) < 3 {
			logger.Fatal("Укажите количество миграций для отката: down <n>")
		}
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			logger.Fatal("Неверное число:", err)
		}
		if err := m.Steps(-n); err != nil && err != migrate.ErrNoChange {
			logger.Fatal("Ошибка отката миграций:", err)
		}
		fmt.Println("Миграции откатаны!")
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			logger.Fatal("Ошибка получения версии:", err)
		}
		fmt.Printf("Версия: %d, Dirty: %t\n", version, dirty)
	case "goto":
		fmt.Println(os.Args)
		if len(os.Args) < 3 {
			logger.Fatal("Укажите версию: goto <version>")
		}
		version, err := strconv.ParseUint(os.Args[2], 10, 64)
		if err != nil {
			logger.Fatal("Неверная версия:", err)
		}
		if err := m.Migrate(uint(version)); err != nil && err != migrate.ErrNoChange {
			logger.Fatal("Ошибка перехода к версии:", err)
		}
		fmt.Println("Переход выполнен!")
	default:
		logger.Fatal("Неизвестная команда:", command)
	}
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
		logger.Fatal("Ошибка создания up-файла:", err)
	}

	// Создание down-файла
	if err := os.WriteFile(downFile, []byte("-- Миграция вниз: "+name+"\n\n"), 0644); err != nil {
		logger.Fatal("Ошибка создания down-файла:", err)
	}

	fmt.Printf("Миграция создана: %s\n", upFile)
	fmt.Printf("Миграция создана: %s\n", downFile)
}
