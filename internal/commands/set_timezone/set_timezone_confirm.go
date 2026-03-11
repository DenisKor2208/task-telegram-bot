// Package set_timezone
package set_timezone

import (
	"strings"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/command"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	btnsettimezoneconfirm "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/set_timezone"
)

// SetTimezoneConfirmCommand - команда для сохранения выбранного часового пояса  "/set_timezone_confirm"
type SetTimezoneConfirmCommand struct{}

func (c *SetTimezoneConfirmCommand) Execute(s command.Model, msg messaging.Message) error {
	tz := strings.TrimSpace(msg.Arguments)
	if tz == "" {
		return s.GetTgClient().SendMessage("Не удалось определить часовой пояс.", msg.UserID)
	}

	// Валидируем, что такой часовой пояс существует
	_, err := time.LoadLocation(tz)
	if err != nil {
		logger.Warn("Пользователь выбрал некорректный часовой пояс", "tz", tz, "user_id", msg.UserID)
		return s.GetTgClient().SendMessage("Некорректный часовой пояс. Пожалуйста, попробуйте ещё раз.", msg.UserID)
	}

	// Получаем пользователя по tg_id
	user, err := s.GetUserStorage().GetUserByTgID(s.GetCtx(), int(msg.UserID))
	if err != nil {
		logger.Error("Не удалось получить пользователя", "user_id", msg.UserID, "error", err)
		return s.GetTgClient().SendMessage("Произошла ошибка. Попробуйте позже.", msg.UserID)
	}

	// Обновляем часовой пояс
	if err := s.GetUserStorage().UpdateUserTimezone(s.GetCtx(), user.ID, tz); err != nil {
		logger.Error("Не удалось обновить часовой пояс", "user_id", user.ID, "error", err)
		return s.GetTgClient().SendMessage("Не удалось сохранить настройки. Попробуйте позже.", msg.UserID)
	}

	// Инициализируем кнопки
	var buttons []bottypes.TgRowButtons

	buttons = append(buttons, btnsettimezoneconfirm.BtnSetTimezoneConfirm...)

	return s.GetTgClient().ShowInlineButtons(resources.TXTSetTimezoneConfirm, buttons, msg.UserID)
}

// NextStep Следующая команда
func (c *SetTimezoneConfirmCommand) NextStep() command.Command {
	return nil
}

// InputField Поле для сохранения данных
func (c *SetTimezoneConfirmCommand) InputField() string {
	return "iana_id"
}
