package records_controller

import (
	"fmt"
	"strings"

	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/helpers"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/model/commands"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/model/currency"
)

type RecordsController struct {
	recordRepo      domain.RecordRepository
	sessionRepo     domain.SessionRepository
	currencyService commands.CurrencyService
	recordsSvc      commands.RecordsService
}

func New(
	s domain.RecordRepository,
	sr domain.SessionRepository,
	cs commands.CurrencyService,
	rs commands.RecordsService) *RecordsController {
	return &RecordsController{
		recordRepo:      s,
		sessionRepo:     sr,
		currencyService: cs,
		recordsSvc:      rs,
	}
}

func (c *RecordsController) AddCommand(msg *domain.Message) *domain.Message {
	line := helpers.TrimLeft(msg.Text, "/add ") //no dup
	record, err := domain.NewRecord(msg.UserID, line)
	if err != nil {
		return msg.Reply(err.Error())
	}

	id, err := c.recordsSvc.Add(msg.Context, msg.UserID, record)
	if err != nil {
		return msg.Reply(err.Error())
	}

	return msg.Reply(fmt.Sprintf("saved with id %d", id))
}

func (c *RecordsController) recordsSum(
	records []domain.Record, cur string) string {
	hash := make(map[string]float64)
	for _, rec := range records {
		hash[rec.Category] += rec.Amount
	}

	var sb strings.Builder
	for cat, val := range hash {
		converted, err := c.currencyService.Convert(
			val, currency.RUB, cur)
		if err == nil {
			s := fmt.Sprintf("%s -- %f %s\n", cat, converted, cur)
			sb.WriteString(s)
		} else {
			logger.Error(err)
		}
	}
	return sb.String()
}

func (c *RecordsController) WeekCommand(msg *domain.Message) *domain.Message {
	records, _ := c.recordRepo.ThisWeek(msg.Context, msg.UserID)
	session, _ := c.sessionRepo.GetOrCreate(msg.Context, msg.UserID)
	return msg.Reply(c.recordsSum(records, session.Currency))
}

func (c *RecordsController) MonthCommand(msg *domain.Message) *domain.Message {
	records, _ := c.recordRepo.ThisMonth(msg.Context, msg.UserID)
	session, _ := c.sessionRepo.GetOrCreate(msg.Context, msg.UserID)
	return msg.Reply(c.recordsSum(records, session.Currency))
}

func (c *RecordsController) YearCommand(msg *domain.Message) *domain.Message {
	records, _ := c.recordRepo.ThisYear(msg.Context, msg.UserID)
	session, _ := c.sessionRepo.GetOrCreate(msg.Context, msg.UserID)
	return msg.Reply(c.recordsSum(records, session.Currency))
}

func (c *RecordsController) SelectCurrencyCommand(msg *domain.Message) *domain.Message {
	session, _ := c.sessionRepo.GetOrCreate(msg.Context, msg.UserID)

	buttons := make(map[string]string)
	for _, currency := range c.currencyService.GetValid() {
		buttons[currency] = "/set_currency " + currency
	}

	return &domain.Message{
		Text:    "selected " + session.Currency + " select new",
		UserID:  msg.UserID,
		Buttons: buttons,
	}
}

func (c *RecordsController) SetCurrencyCommand(msg *domain.Message) *domain.Message {
	line := helpers.TrimLeft(msg.Text, "/set_currency ")
	reader := strings.NewReader(line)
	var cur string
	fmt.Fscanf(reader, "%s", &cur)

	if !c.currencyService.IsValid(cur) {
		return msg.Reply("Invalid currency")
	}

	session, _ := c.sessionRepo.GetOrCreate(msg.Context, msg.UserID)
	session.Currency = cur
	_ = c.sessionRepo.Save(msg.Context, session)

	return msg.Reply("you selected " + cur)
}

func (c *RecordsController) LimitCommand(msg *domain.Message) *domain.Message {
	line := helpers.TrimLeft(msg.Text, "/limit ")
	reader := strings.NewReader(line)
	var limit float64
	fmt.Fscanf(reader, "%f", &limit)

	session, _ := c.sessionRepo.GetOrCreate(msg.Context, msg.UserID)
	session.MonthLimit = limit
	logger.Debug("saving new limit", limit)
	err := c.sessionRepo.Save(msg.Context, session)

	if err != nil {
		return msg.Reply(err.Error())
	} else {
		return msg.Reply(
			fmt.Sprintf("set limit %f %s", limit, session.Currency),
		)
	}
}
