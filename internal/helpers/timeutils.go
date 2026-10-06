package helpers

import (
	"context"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/storage"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
)

const (
	// DeadlineLayout - формат отображения дедлайна пользователю.
	DeadlineLayout = "02.01.2006 15:04"

	// DeadlineInputLayout - формат разбора введённого дедлайна:
	// день, месяц и час — 1 или 2 цифры, год — 4 цифры, минуты — 2 цифры.
	DeadlineInputLayout = "2.1.2006 15:04"

	// DefaultDeadlineTime - время дедлайна, если пользователь указал только дату (до конца дня).
	DefaultDeadlineTime = "23:59"
)

// UserLocation возвращает часовой пояс пользователя с Telegram ID tgID.
// Если пользователь не найден или его часовой пояс некорректен, возвращает UTC
// и пишет предупреждение в лог — чтобы список задач всё равно можно было показать.
func UserLocation(ctx context.Context, userStorage storage.User, tgID int64) *time.Location {
	user, err := userStorage.GetUserByTgID(ctx, int(tgID))
	if err != nil {
		logger.Warn("Не удалось получить пользователя, используется UTC", "user_id", tgID, "error", err)
		return time.UTC
	}

	loc, err := time.LoadLocation(user.Timezone)
	if err != nil {
		logger.Warn("Не удалось загрузить часовой пояс, используется UTC",
			"user_id", tgID, "timezone", user.Timezone, "error", err)
		return time.UTC
	}

	return loc
}

// DeadlineLabel возвращает подпись дедлайна для кнопки задачи в часовом поясе loc:
// «до 10.10.2026 18:00» или «без дедлайна».
func DeadlineLabel(deadline *time.Time, loc *time.Location) string {
	if deadline == nil {
		return "без дедлайна"
	}
	return "до " + deadline.In(loc).Format(DeadlineLayout)
}

// DeadlineText возвращает строку о дедлайне для сообщений в часовом поясе loc:
// «Дедлайн: 10.10.2026 18:00» или «Без дедлайна».
func DeadlineText(deadline *time.Time, loc *time.Location) string {
	if deadline == nil {
		return "Без дедлайна"
	}
	return "Дедлайн: " + deadline.In(loc).Format(DeadlineLayout)
}
