// Package helpers
package helpers

import (
	"fmt"
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

// TaskInput — результат разбора текста новой задачи.
type TaskInput struct {
	Description     string // описание задачи (пусто, если введена только дата)
	Date            string // дата дедлайна "Д.М.ГГГГ" (пусто — без дедлайна)
	Time            string // время дедлайна "Ч:ММ" (пусто — время не указано)
	TimeWithoutDate bool   // в конце текста время "ЧЧ:ММ", но перед ним нет даты
}

var (
	// <описание> <дата> <время>: день и месяц — 1–2 цифры, год — 4 цифры,
	// время через двоеточие или точку (18:00 или 18.00). Описание может отсутствовать.
	taskDateTimeRe = regexp.MustCompile(`^(?:(.*?)\s+)?(\d{1,2}\.\d{1,2}\.\d{4})\s+(\d{1,2})[:.](\d{2})$`)
	// <описание> <дата> — дедлайн без времени. Описание может отсутствовать.
	taskDateRe = regexp.MustCompile(`^(?:(.*?)\s+)?(\d{1,2}\.\d{1,2}\.\d{4})$`)
	// В конце текста время без даты. Только через двоеточие: время через точку
	// без даты не отличить от обычного числа в описании (например, «Купить 1.50 кг»).
	taskTimeOnlyRe = regexp.MustCompile(`(?:^|\s)\d{1,2}:\d{2}$`)
)

// ParseTaskInput разбирает текст новой задачи: "<описание> [<дата> [<время>]]".
// Возвращает false, если текст пустой.
func ParseTaskInput(input string) (TaskInput, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return TaskInput{}, false
	}

	if m := taskDateTimeRe.FindStringSubmatch(trimmed); m != nil {
		return TaskInput{
			Description: strings.TrimSpace(m[1]),
			Date:        m[2],
			Time:        m[3] + ":" + m[4],
		}, true
	}

	if m := taskDateRe.FindStringSubmatch(trimmed); m != nil {
		return TaskInput{
			Description: strings.TrimSpace(m[1]),
			Date:        m[2],
		}, true
	}

	// Если не дата, считаем описанием
	return TaskInput{
		Description:     trimmed,
		TimeWithoutDate: taskTimeOnlyRe.MatchString(trimmed),
	}, true
}

func GetDisplayName(msg messaging.Message) string {
	if msg.UserDisplayName != "" {
		return msg.UserDisplayName
	}
	return msg.UserName
}

// GetUserNameForDB возвращает имя пользователя для сохранения в БД.
// В БД стоит проверка, что имя не пустое, а @username в Telegram необязателен,
// поэтому берётся первое непустое значение: @username, имя из профиля, user_<Telegram ID>.
func GetUserNameForDB(msg messaging.Message) string {
	if msg.UserName != "" {
		return msg.UserName
	}
	if msg.UserDisplayName != "" {
		return msg.UserDisplayName // имя и фамилия из профиля Telegram
	}
	return fmt.Sprintf("user_%d", msg.UserID)
}
