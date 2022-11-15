package tg

import (
	"context"
	"runtime"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
	messages "gitlab.ozon.dev/albatros2002/telegram-bot/internal/model/messages"
)

type TokenGetter interface {
	Token() string
}

type Client struct {
	client *tgbotapi.BotAPI
	ctx    context.Context
}

func New(tokenGetter TokenGetter, ctx context.Context) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}
	return &Client{
		client: client,
		ctx:    ctx,
	}, nil
}

func makeKeyboard(buttons map[string]string) tgbotapi.InlineKeyboardMarkup {
	keyboardButtons := []tgbotapi.InlineKeyboardButton{}

	for label, data := range buttons {
		b := tgbotapi.NewInlineKeyboardButtonData(label, data)
		keyboardButtons = append(keyboardButtons, b)
	}
	row := tgbotapi.NewInlineKeyboardRow(keyboardButtons...)
	return tgbotapi.NewInlineKeyboardMarkup(row)
}

func (c *Client) SendMessage(m domain.Message) error {
	msg := tgbotapi.NewMessage(m.UserID, m.Text)
	if m.Buttons != nil {
		msg.ReplyMarkup = makeKeyboard(m.Buttons)
	} else {
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}
	_, err := c.client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}
	return nil
}

func (c *Client) ListenUpdates(msgModel *messages.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	logger.Info("TGClient listening for messages")

	for {
		select {
		case update := <-updates:

			if update.CallbackQuery != nil {
				logger.Debug("callback query")
				logger.Debugf("%s", update.CallbackQuery.From.ID)
				logger.Debugf("%+v", update.CallbackQuery.Data)
				err := msgModel.IncomingMessage(domain.Message{
					Text:    update.CallbackQuery.Data,
					UserID:  update.CallbackQuery.From.ID,
					Context: context.Background(),
				})
				if err != nil {
					logger.Warn("Error processing message:", err)
				}
			}

			if update.Message != nil {
				logger.Debugf("%s %s", update.Message.From.UserName, update.Message.Text)

				err := msgModel.IncomingMessage(domain.Message{
					Text:    update.Message.Text,
					UserID:  update.Message.From.ID,
					Context: context.Background(),
				})
				if err != nil {
					logger.Warn("Error processing message:", err)
				}
			}

		case <-c.ctx.Done():
			logger.Debug("TGClient ctx done")
			c.client.StopReceivingUpdates()
			n := runtime.NumGoroutine()
			logger.Debug("Number of goroutines: ", n)
			logger.Debug("TGClient done")
			return
		}
	}
}
