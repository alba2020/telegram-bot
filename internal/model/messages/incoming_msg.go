package messages

import (
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
)

type MessageSender interface {
	SendMessage(domain.Message) error
}

type Model struct {
	tgClient MessageSender
	router   domain.Router
}

func New(tgClient MessageSender, router domain.Router) *Model {
	return &Model{
		tgClient: tgClient,
		router:   router,
	}
}

func (m *Model) send(msg domain.Message) {
	err := m.tgClient.SendMessage(msg)
	if err != nil {
		logger.Error("could not send message", err.Error())
	}
}

func (m *Model) IncomingMessage(msg domain.Message) error {
	handler, err := m.router.Route(&msg)

	if err != nil {
		return m.tgClient.SendMessage(domain.Message{
			Text:    err.Error(),
			UserID:  msg.UserID,
			Buttons: nil,
		})
	}

	response := handler(&msg)
	if response.Text != "" {
		m.send(*response)
	} else {
		logger.Warn("response is empty")
	}

	return nil
}
