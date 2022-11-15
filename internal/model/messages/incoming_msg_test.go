package messages

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/ctrl/info_controller"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/ctrl/records_controller"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	mocks "gitlab.ozon.dev/albatros2002/telegram-bot/internal/mocks/messages"
	currencyfake "gitlab.ozon.dev/albatros2002/telegram-bot/internal/model/currency/currency_fake"
	recordmemoryrepo "gitlab.ozon.dev/albatros2002/telegram-bot/internal/repo/record_memory_repo"
	sessionmemoryrepo "gitlab.ozon.dev/albatros2002/telegram-bot/internal/repo/session_memory_repo"
	recordsservice "gitlab.ozon.dev/albatros2002/telegram-bot/internal/services/records_service"
)

const userId int64 = 132

func getRouter() domain.Router {
	recordRepo := recordmemoryrepo.New()
	sessionRepo, _ := sessionmemoryrepo.New()
	currencyService, _ := currencyfake.New()

	var infoController = info_controller.New()

	recordsSvc, _ := recordsservice.New(
		sessionRepo, recordRepo, currencyService)

	var recordsController = records_controller.New(
		recordRepo, sessionRepo, currencyService, recordsSvc)

	var router = domain.NewRouter(map[string]domain.CommandHandler{
		"/start":           infoController.StartCommand,
		"/ping":            infoController.PingCommand,
		"/add":             recordsController.AddCommand,
		"/week":            recordsController.WeekCommand,
		"/month":           recordsController.MonthCommand,
		"/year":            recordsController.YearCommand,
		"/select_currency": recordsController.SelectCurrencyCommand,
		"/set_currency":    recordsController.SetCurrencyCommand,
	})
	return router
}

var router = getRouter()

func Test_OnStartCommand_ShouldAnswerWithIntroMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender, router)

	sender.EXPECT().SendMessage(domain.Message{
		Text:   "hello",
		UserID: userId,
	})

	err := model.IncomingMessage(domain.Message{
		Text:   "/start",
		UserID: userId,
	})
	time.Sleep(time.Second * 1)
	assert.NoError(t, err)
}

func Test_OnUnknownCommand_ShouldAnswerWithErrorMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender, router)

	sender.EXPECT().SendMessage(domain.Message{
		Text:   "handler not found",
		UserID: userId,
	})

	err := model.IncomingMessage(domain.Message{
		Text:   "/qwerty",
		UserID: userId,
	})
	time.Sleep(time.Second * 1)
	assert.NoError(t, err)
}

func Test_OnAddCommand_ShouldAnswerWithSavedId(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	model := New(sender, router)

	sender.EXPECT().SendMessage(domain.Message{
		Text:    "saved with id 1",
		UserID:  userId,
		Buttons: nil,
		Context: nil,
	})

	err := model.IncomingMessage(domain.Message{
		Text:    "/add 1.12 food 15-05-2001",
		UserID:  userId,
		Context: context.Background(),
	})
	time.Sleep(time.Second * 1)
	assert.NoError(t, err)
}
