// Package tg
package tg

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/DenisKor2208/task-telegram-bot/internal/helpers"
	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/messages"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	"github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands"
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

// redactingHTTPClient - HTTP-клиент для запросов к Telegram API, скрывающий токен бота в ошибках.
// Токен входит в адрес каждого запроса (https://api.telegram.org/bot<TOKEN>/...),
// и при сетевой ошибке Go добавляет этот адрес в текст ошибки, откуда он попадает в логи.
type redactingHTTPClient struct {
	client *http.Client
	token  string
}

// Do выполняет запрос и заменяет токен в адресе сетевой ошибки на "<TOKEN>".
func (c *redactingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)

	var urlErr *url.Error
	if err != nil && c.token != "" && errors.As(err, &urlErr) {
		urlErr.URL = strings.ReplaceAll(urlErr.URL, c.token, "<TOKEN>")
	}

	return resp, err
}

func New(tokenGetter TokenGetter, handlerProcessingFunc HandlerFunc) (*Client, error) {
	token := tokenGetter.Token()

	// Все запросы библиотеки идут через этот клиент, поэтому токен не попадёт в логи
	// ни из нашего кода, ни из внутренних логов библиотеки.
	httpClient := &redactingHTTPClient{client: &http.Client{}, token: token}

	client, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, httpClient)
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
	deleteMsg := func(chatID int64, msgID int) {
		if err := deleteInlineButtons(c, chatID, msgID); err != nil {
			logger.Error("Ошибка удаления кнопок:", "err", err)
		}
		if err := deleteMessage(c, chatID, msgID); err != nil {
			logger.Error("Ошибка удаления сообщения:", "err", err)
		}
	}

	parseMessage := func(text string, msg *messaging.Message) {
		if strings.HasPrefix(text, "/") {
			cmd, args := helpers.ParseCommandAndArgs(text)
			msg.IsCommand = true
			msg.Command = cmd
			msg.Arguments = args
		}
	}

	if tgUpdate.Message != nil {
		// Пользователь написал текстовое сообщение.
		logger.Info(fmt.Sprintf("[%s][%v] %s", tgUpdate.Message.From.UserName, tgUpdate.Message.From.ID, tgUpdate.Message.Text))

		msg := messaging.Message{
			Text:            tgUpdate.Message.Text,
			UserID:          tgUpdate.Message.From.ID,
			UserName:        tgUpdate.Message.From.UserName,
			UserDisplayName: strings.TrimSpace(tgUpdate.Message.From.FirstName + " " + tgUpdate.Message.From.LastName),
			Date:            int64(tgUpdate.Message.Date),
			/**/
			Arguments: strings.TrimSpace(tgUpdate.Message.Text),
		}

		// Удаление кнопок и сообщения
		deleteMsg(msg.UserID, tgUpdate.Message.MessageID)

		// Определяем, команда ли это
		parseMessage(msg.Text, &msg)

		if err := msgModel.IncomingMessage(msg); err != nil {
			showProcessingError(c, msg.UserID, err)

			logger.Error("error processing message:", "err", err)
		}
	} else if tgUpdate.CallbackQuery != nil {
		// Пользователь нажал кнопку.
		logger.Info(fmt.Sprintf("[%s][%v] Callback: %s", tgUpdate.CallbackQuery.From.UserName, tgUpdate.CallbackQuery.From.ID, tgUpdate.CallbackQuery.Data))

		callback := tgbotapi.NewCallback(tgUpdate.CallbackQuery.ID, tgUpdate.CallbackQuery.Data)
		if _, err := c.Client.Request(callback); err != nil {
			logger.Error("Ошибка Request callback:", "err", err)
		}

		// Удаление кнопок и сообщения
		deleteMsg(tgUpdate.CallbackQuery.From.ID, tgUpdate.CallbackQuery.Message.MessageID)

		msg := messaging.Message{
			Text:            tgUpdate.CallbackQuery.Data,
			UserID:          tgUpdate.CallbackQuery.From.ID,
			UserName:        tgUpdate.CallbackQuery.From.UserName,
			UserDisplayName: strings.TrimSpace(tgUpdate.CallbackQuery.From.FirstName + " " + tgUpdate.CallbackQuery.From.LastName),
			IsCallback:      true,
			CallbackMsgID:   tgUpdate.CallbackQuery.InlineMessageID,
			Date:            int64(tgUpdate.CallbackQuery.Message.Date),
		}

		// Определяем, команда ли это
		parseMessage(msg.Text, &msg)

		if err := msgModel.IncomingMessage(msg); err != nil {
			showProcessingError(c, msg.UserID, err)

			logger.Error("error processing callback:", "err", err)
		}
	}
}

// showProcessingError показывает пользователю сообщение об ошибке обработки.
// Причину ошибки показывает, только если пользователь может исправить её сам (messaging.UserError),
// иначе — общее сообщение, чтобы не раскрывать внутренние детали (подробности пишутся в лог).
func showProcessingError(c *Client, userID int64, err error) {
	text := resources.ErrGeneric

	var userErr *messaging.UserError
	if errors.As(err, &userErr) {
		text = userErr.Text
	}

	_ = c.ShowInlineButtons(text, commands.BtnOther, userID)
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

	return nil
}
