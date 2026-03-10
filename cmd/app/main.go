package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
)

func main() {
	logger.Info("Старт приложения")

	ctx := context.Background()

	// Контекст для graceful shutdown
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGHUP,  //"Hang up" — обычно отправляется при закрытии терминала.
		syscall.SIGINT,  //"Interrupt" — отправляется при нажатии Ctrl+C в терминале.
		syscall.SIGTERM, //"Terminate" — стандартный сигнал для завершения процесса (от systemd, kill).
		syscall.SIGQUIT, //"Quit" — аналогично SIGINT, но может вызывать дамп стека (stack trace).
	)
	defer cancel()

	// Создаём приложение
	app, err := NewApp(ctx)
	if err != nil {
		logger.Fatal("Ошибка инициализации приложения:", "err", err)
	}

	// Логируем успешное подключение
	log.Printf("Authorized on account %s", app.tgClient.Client.Self.UserName)

	// Запускаем Telegram-клиент
	go app.tgClient.ListenUpdates(app.msgModel)

	// Запускаем проверку просроченных задач
	go app.overdueChecker.Start(ctx)

	// Ждём сигнала завершения
	<-ctx.Done()

	// Graceful shutdown: закрываем Redis
	if app.redisClient != nil {
		app.redisClient.Close()
	}

	logger.Info("Application shutdown complete")
}
