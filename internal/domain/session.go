package domain

import (
	"context"

	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/model/currency"
)

type Session struct {
	UserId     int64   `db:"user_id"`
	Currency   string  `db:"currency"`
	MonthLimit float64 `db:"month_limit"`
}

type SessionRepository interface {
	Find(ctx context.Context, userId int64) (*Session, error)
	GetOrCreate(ctx context.Context, userId int64) (*Session, error)
	Save(context.Context, *Session) error
}

func NewSession(userId int64) (*Session, error) {
	return &Session{
		UserId:     userId,
		Currency:   currency.RUB,
		MonthLimit: -1,
	}, nil
}
