package commands

import (
	"context"

	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
)

type RecordsService interface {
	Add(context.Context, int64, *domain.Record) (int64, error)
}

type CurrencyService interface {
	IsValid(string) bool
	GetValid() []string
	Convert(float64, string, string) (float64, error)
}
