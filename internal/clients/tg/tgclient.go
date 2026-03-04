package tg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"

	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/bottypes"
)

type HandlerFunc func(tgUpdate *tgbotapi.Update, c *Client, msgModel *messages.Model)

func (f HandlerFunc) RunFunc(tgUpdate *tgbotapi.Update, c *Client, msgModel *messages.Model) {
	f(tgUpdate, c, msgModel)
}

type Client struct {
	Client                *tgbotapi.BotAPI // Клиент Telegram API
	handlerProcessingFunc HandlerFunc      // Функция обработки входящих сообщений
}

type TokenGetter interface {
	Token() string
}

func New(tokenGetter TokenGetter, handlerProcessingFunc HandlerFunc) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка NewBotAPI")
	}

	client.Debug = true

	return &Client{
		Client:                client,
		handlerProcessingFunc: handlerProcessingFunc,
	}, nil
}

func (c *Client) SendMessage(text string, userID int64) error {
	msg := tgbotapi.NewMessage(userID, text)
	msg.ParseMode = "HTML"
	_, err := c.Client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "Ошибка отправки сообщения client.Send")
	}
	return nil
}

func (c *Client) ListenUpdates(msgModel *messages.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.Client.GetUpdatesChan(u)

	logger.Info("Start listening for tg messages")

	for {
		select {
		case update, ok := <-updates:
			if !ok {
				// Канал закрыт — выходим (редкий случай, но для надёжности).
				logger.Info("Updates channel closed")
				return
			}
			// Функция обработки сообщений (обернутая в middleware).
			c.handlerProcessingFunc.RunFunc(&update, c, msgModel)
			// Вместо ProcessingMessages(update, c, msgModel)
		case <-msgModel.GetCtx().Done():
			// Контекст отменён (по сигналу, например Ctrl+C) — выходим gracefully.
			logger.Info("Context cancelled, stopping listener")
			return
		}
	}
}

// ProcessingMessages функция обработки сообщений.
func ProcessingMessages(tgUpdate *tgbotapi.Update, c *Client, msgModel *messages.Model) {
	if tgUpdate.Message != nil {
		// Пользователь написал текстовое сообщение.
		logger.Info(fmt.Sprintf("[%s][%v] %s", tgUpdate.Message.From.UserName, tgUpdate.Message.From.ID, tgUpdate.Message.Text))

		isCommand := false
		command, arguments := parseCommandAndArgs(tgUpdate.Message.Text)
		if command != "" {
			isCommand = true
		}

		err := msgModel.IncomingMessage(messaging.Message{
			Text:            tgUpdate.Message.Text,
			Command:         command,
			Arguments:       arguments,
			IsCommand:       isCommand,
			UserID:          tgUpdate.Message.From.ID,
			UserName:        tgUpdate.Message.From.UserName,
			UserDisplayName: strings.TrimSpace(tgUpdate.Message.From.FirstName + " " + tgUpdate.Message.From.LastName),
			Date:            int64(tgUpdate.Message.Date),
		})
		if err != nil {
			logger.Error("error processing message:", "err", err)
		}
	} else if tgUpdate.CallbackQuery != nil {
		// Пользователь нажал кнопку.
		logger.Info(fmt.Sprintf("[%s][%v] Callback: %s", tgUpdate.CallbackQuery.From.UserName, tgUpdate.CallbackQuery.From.ID, tgUpdate.CallbackQuery.Data))
		callback := tgbotapi.NewCallback(tgUpdate.CallbackQuery.ID, tgUpdate.CallbackQuery.Data)
		if _, err := c.Client.Request(callback); err != nil {
			logger.Error("Ошибка Request callback:", "err", err)
		}

		if err := deleteInlineButtons(c, tgUpdate.CallbackQuery.From.ID, tgUpdate.CallbackQuery.Message.MessageID); err != nil {
			logger.Error("Ошибка удаления кнопок:", "err", err)
		}

		if err := deleteMessage(c, tgUpdate.CallbackQuery.From.ID, tgUpdate.CallbackQuery.Message.MessageID); err != nil {
			logger.Error("Ошибка удаления кнопок:", "err", err)
		}

		// buttonText := tgUpdate.CallbackQuery.Data
		// parts := strings.Fields(strings.TrimPrefix(buttonText, "/"))
		// command, arguments, isCommand := "", "", strings.HasPrefix(buttonText, "/")
		// if isCommand && len(parts) > 0 {
		// 	command = parts[0]
		// 	arguments = strings.Join(parts[1:], " ")
		// }

		isCommand := false
		command, arguments := parseCommandAndArgs(tgUpdate.CallbackQuery.Data)
		if command != "" {
			isCommand = true
		}

		err := msgModel.IncomingMessage(messaging.Message{
			Text:            tgUpdate.CallbackQuery.Data,
			Command:         command,
			Arguments:       arguments,
			IsCommand:       isCommand,
			UserID:          tgUpdate.CallbackQuery.From.ID,
			UserName:        tgUpdate.CallbackQuery.From.UserName,
			UserDisplayName: strings.TrimSpace(tgUpdate.CallbackQuery.From.FirstName + " " + tgUpdate.CallbackQuery.From.LastName),
			IsCallback:      true,
			CallbackMsgID:   tgUpdate.CallbackQuery.InlineMessageID,
			Date:            int64(tgUpdate.CallbackQuery.Message.Date),
		})
		if err != nil {
			logger.Error("error processing message from callback:", "err", err)
		}
	}
}

// ShowInlineButtons Отображение кнопок меню под сообщением с ответом.
// Их нажатие ожидает коллбек-ответ.
func (c *Client) ShowInlineButtons(text string, buttons []bottypes.TgRowButtons, userID int64) error {
	keyboard := make([][]tgbotapi.InlineKeyboardButton, len(buttons))
	for i := 0; i < len(buttons); i++ {
		tgRowButtons := buttons[i]
		keyboard[i] = make([]tgbotapi.InlineKeyboardButton, len(tgRowButtons))
		for j := 0; j < len(tgRowButtons); j++ {
			tgInlineButton := tgRowButtons[j]
			keyboard[i][j] = tgbotapi.NewInlineKeyboardButtonData(tgInlineButton.DisplayName, tgInlineButton.Value)
		}
	}
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	msg := tgbotapi.NewMessage(userID, text)
	msg.ReplyMarkup = numericKeyboard
	msg.ParseMode = "HTML"
	_, err := c.Client.Send(msg)
	if err != nil {
		logger.Error("Ошибка отправки сообщения", "err", err)
		return errors.Wrap(err, "client.Send with inline-buttons")
	}
	return nil
}

// SendMessageWithTemplate Отправка сообщения с введенным шаблоном в поле ввода.
func (c *Client) SendMessageWithTemplate(text string, textTemplate string, userID int64) error {
	msg := tgbotapi.NewMessage(userID, text)
	msg.ParseMode = "HTML"

	// Включаем force_reply — поднимает клавиатуру с ответом
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: true,
		// Опционально: можно добавить placeholder
		InputFieldPlaceholder: textTemplate,
		Selective:             false,
	}

	_, err := c.Client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "Ошибка отправки сообщения client.Send")
	}
	return nil
}

func deleteInlineButtons(c *Client, userID int64, msgID int) error {
	// Создаём пустую inline-клавиатуру
	emptyKeyboard := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
	}

	// Удаляем inline-кнопки
	msg := tgbotapi.NewEditMessageReplyMarkup(userID, msgID, emptyKeyboard)
	_, err := c.Client.Send(msg)
	if err != nil {
		logger.Error("Ошибка удаления inline-кнопок", "err", err)
		return errors.Wrap(err, "client.Send delete inline-buttons")
	}
	return nil
}

/* func deleteMessage(c *Client, chatID int64, msgID int) error {
	deleteConfig := tgbotapi.NewDeleteMessage(chatID, msgID)
	_, err := c.Client.Send(deleteConfig)
	if err != nil {
		// Обработка ошибок (например, сообщение слишком старое или уже удалено)
		logger.Error("Ошибка удаления сообщения", "err", err)
		return errors.Wrap(err, "deleteMessage failed")
	}
	return nil
} */

func deleteMessage(c *Client, chatID int64, msgID int) error {
	deleteConfig := tgbotapi.NewDeleteMessage(chatID, msgID)

	// Используем Request вместо Send
	resp, err := c.Client.Request(deleteConfig)
	if err != nil {
		logger.Error("Ошибка запроса на удаление сообщения", "err", err)
		return errors.Wrap(err, "deleteMessage request failed")
	}

	// Проверяем, что API вернуло успешный статус
	if !resp.Ok {
		logger.Error("Telegram API вернул ошибку при удалении",
			"error_code", resp.ErrorCode,
			"description", resp.Description)
		return errors.Errorf("deleteMessage failed: %s", resp.Description)
	}

	// Можно добавить проверку, что в результате пришло true,
	// но обычно если Ok == true, то всё хорошо.
	// Если нужно явно проверить:
	// if resp.Result != true { ... }

	return nil
}

func parseCommandAndArgs(input string) (command string, args string) {
	trimmedInput := strings.TrimSpace(input)

	if trimmedInput == "" {
		return "", ""
	}

	//re := regexp.MustCompile(`^/([^\s]+)(?:\s+(.*))?`) // старый
	re := regexp.MustCompile(`^/([^\s]+)\s*(.*)`) // новый
	matches := re.FindStringSubmatch(trimmedInput)

	if len(matches) > 0 {
		command = strings.TrimSpace(matches[1])
		args = strings.TrimSpace(matches[2])
	} else {
		command = ""
		args = trimmedInput
	}

	return command, args
}
