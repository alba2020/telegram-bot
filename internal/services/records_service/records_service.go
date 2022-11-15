package recordsservice

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/model/commands"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/model/currency"
)

type RecordsService struct {
	sessionRepo domain.SessionRepository
	recordRepo  domain.RecordRepository
	currencySvc commands.CurrencyService
}

func New(
	sessionRepo domain.SessionRepository,
	recordRepo domain.RecordRepository,
	currencySvc commands.CurrencyService) (*RecordsService, error) {
	return &RecordsService{
		sessionRepo: sessionRepo,
		recordRepo:  recordRepo,
		currencySvc: currencySvc,
	}, nil
}

func (svc *RecordsService) Add(ctx context.Context, userId int64, rec *domain.Record) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "records_service: add")
	defer span.Finish()

	session, _ := svc.sessionRepo.GetOrCreate(ctx, userId)

	rec.Amount, _ = svc.currencySvc.Convert(
		rec.Amount, session.Currency, currency.RUB)
	logger.Debug("Saving amount", rec.Amount)
	limit, _ := svc.currencySvc.Convert(
		session.MonthLimit, session.Currency, currency.RUB)
	logger.Debug("This month limit", limit)
	spent, _ := svc.recordRepo.ThisMonthSum(ctx, userId)
	logger.Debug("This month spent", spent)

	if limit > 0 && (spent+rec.Amount) > limit {
		return -1, errors.New("limit exceeded")
	}

	id, err := svc.recordRepo.Save(ctx, rec)
	if err != nil {
		return 0, err
	}

	return id, nil
}
