package dbutils

import (
	"bytes"   // Для работы с буфером байтов — помогает строить строки эффективно без лишних аллокаций
	"context" // Для передачи контекста в метод Log
	"fmt"     // Для форматирования значений в строке
	"strings" // Для замены пробелов на подчёркивания в логах
	"time"    // Для работы с длительностью времени в настройках соединения

	"github.com/DenisKor2208/task-telegram-bot/internal/logger" // Кастомный логгер для вывода сообщений
	"github.com/jackc/pgx/v5"                                   // Драйвер для PostgreSQL — предоставляет конфиг и логгер
	"github.com/jackc/pgx/v5/stdlib"                            // Стандартная библиотека pgx для интеграции с database/sql
	"github.com/jackc/pgx/v5/tracelog"                          // Пакет tracelog для работы с Tracer и логированием
	"github.com/jmoiron/sqlx"                                   // Расширение database/sql для удобной работы с SQL (поддержка структур)
)

// pgxLogger реализует интерфейс tracelog.Logger для перенаправления логов в кастомный логгер.
// Это позволяет перехватывать логи от pgx и выводить их через наш logger с кастомным форматированием.
type pgxLogger struct{}

// Log — метод, который вызывается tracelog для логирования сообщений.
// Он строит полное сообщение (основное + данные), форматирует его и выводит в logger по уровню.
func (l *pgxLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	// Создаём буфер для эффективного построения строки лога без лишних копирований
	var buf bytes.Buffer
	buf.WriteString(msg) // Добавляем основное сообщение

	// Если есть дополнительные данные (параметры), добавляем их в лог для отладки
	if len(data) > 0 {
		buf.WriteString(" | data: ") // Префикс для разделения основного сообщения и данных
		for k, v := range data {     // Проходим по каждой паре ключ-значение в data
			buf.WriteString(k)   // Добавляем ключ
			buf.WriteString("=") // Разделитель между ключом и значением
			// Форматируем значение как строку, заменяем пробелы на _ для лучшей читаемости в логах
			buf.WriteString(strings.ReplaceAll(fmt.Sprintf("%v", v), " ", "_"))
			buf.WriteString(" ") // Пробел между парами для аккуратного вида
		}
		// Удаляем лишний пробел в конце, если данные были добавлены
		logMsg := strings.TrimSuffix(buf.String(), " ")
		buf.Reset()
		buf.WriteString(logMsg)
	}

	// Получаем итоговую строку из буфера
	logMsg := buf.String()

	// Выбираем уровень логирования и выводим сообщение в наш кастомный logger
	switch level {
	case tracelog.LogLevelTrace, tracelog.LogLevelDebug: // Для детальной отладки
		logger.Debug(logMsg)
	case tracelog.LogLevelInfo: // Для информационных сообщений
		logger.Info(logMsg)
	case tracelog.LogLevelWarn: // Для предупреждений
		logger.Warn(logMsg)
	case tracelog.LogLevelError: // Для ошибок
		logger.Error(logMsg)
	default: // По умолчанию — как info, если уровень неизвестен
		logger.Info(logMsg)
	}
}

// NewDBConnect устанавливает соединение с PostgreSQL через sqlx и pgx.
// Принимает строку подключения (DSN) и возвращает готовый *sqlx.DB или ошибку.
func NewDBConnect(connString string) (*sqlx.DB, error) {
	// Парсим строку подключения в конфиг pgx — это позволяет настроить параметры соединения
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		// Если парсинг не удался, логируем ошибку и возвращаем её
		logger.Error("Ошибка парсинга строки подключения: " + err.Error())
		return nil, err
	}

	// Устанавливаем имя приложения для идентификации в БД (полезно для мониторинга)
	config.RuntimeParams["application_name"] = "tg-bot"

	// Создаём TraceLogTracer с нашим кастомным логгером и уровнем логирования
	tracer := &tracelog.TraceLog{
		Logger:   &pgxLogger{},           // Кастомный логгер
		LogLevel: tracelog.LogLevelDebug, // Уровень логирования (debug для подробностей)
	}
	// Подключаем Tracer к конфигу
	config.Tracer = tracer

	// Регистрируем конфиг в stdlib для совместимости с database/sql (нужно для sqlx)
	connStr := stdlib.RegisterConnConfig(config)

	// Подключаемся к БД через sqlx, используя драйвер pgx
	dbh, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		// Если подключение не удалось, логируем и возвращаем ошибку
		logger.Error("Ошибка подключения к БД: " + err.Error())
		return nil, err
	}

	// Настройки пула соединений для оптимизации производительности
	// Эти параметры контролируют, сколько соединений держать открытыми и как долго
	dbh.SetMaxOpenConns(10)                 // Максимум открытых соединений одновременно (предотвращает перегрузку БД)
	dbh.SetMaxIdleConns(5)                  // Максимум "праздных" (неиспользуемых) соединений в пуле
	dbh.SetConnMaxLifetime(5 * time.Minute) // Максимальное время жизни одного соединения (обновляет их периодически)

	// Возвращаем готовый объект sqlx.DB для использования в коде
	return dbh, nil
}
