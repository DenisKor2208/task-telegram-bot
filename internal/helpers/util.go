// Package helpers
package helpers

import (
	"regexp"
	"strings"

	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
)

// ParseCommandAndArgs парсит входную строку и извлекает команду (начинающуюся с "/") и аргументы.
// Возвращает команду и аргументы.
func ParseCommandAndArgs(input string) (command string, args string) {
	trimmedInput := strings.TrimSpace(input)

	if trimmedInput == "" {
		return "", ""
	}

	commandRegex := regexp.MustCompile(`^/([^\s]+)\s*(.*)`)

	matches := commandRegex.FindStringSubmatch(trimmedInput)
	if len(matches) > 0 {
		command = /* "/" +  */ strings.TrimSpace(matches[1])
		args = strings.TrimSpace(matches[2])
	} else {
		command = ""
		args = trimmedInput
	}

	return command, args
}

// ParseForCommandSaveTask парсит строку в формате "/save_task <описание> <дата в DD.MM.YYYY> <время в HH:MM>".
// Возвращает описание, дату и время в формате "DD.MM.YYYY HH:MM", и флаг успешности.
func ParseForCommandSaveTask(input string) (description, datetime string, ok bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", "", false
	}

	fullPattern := `^(.+)\s+(\d{2}\.\d{2}\.\d{4})\s+(\d{2}:\d{2})$`
	fullRe := regexp.MustCompile(fullPattern)
	if matches := fullRe.FindStringSubmatch(trimmed); len(matches) == 4 {
		description = strings.TrimSpace(matches[1])
		datetime = matches[2] + " " + matches[3]
		return description, datetime, true
	}

	datePattern := `^\d{2}\.\d{2}\.\d{4}\s+\d{2}:\d{2}$`
	dateRe := regexp.MustCompile(datePattern)
	if dateRe.MatchString(trimmed) {
		return "", trimmed, true
	}

	// Если не пусто и не дата, считаем описанием
	return trimmed, "", true
}

func GetDisplayName(msg messaging.Message) string {
	if msg.UserDisplayName != "" {
		return msg.UserDisplayName
	}
	return msg.UserName
}
