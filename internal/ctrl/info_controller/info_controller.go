package info_controller

import (
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
)

type InfoController struct{}

func New() *InfoController {
	return &InfoController{}
}

func (c *InfoController) StartCommand(msg *domain.Message) *domain.Message {
	return msg.Reply("hello")
}

func (c *InfoController) PingCommand(msg *domain.Message) *domain.Message {
	span, _ := opentracing.StartSpanFromContext(msg.Context, "ping")
	defer span.Finish()

	return msg.Reply("pong!")
}
