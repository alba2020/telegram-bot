package domain

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

type Record struct {
	Id       int64     `db:"id"`
	UserId   int64     `db:"user_id"`
	Amount   float64   `db:"amount"`
	Category string    `db:"category"`
	Date     time.Time `db:"date"`
}

type RecordRepository interface {
	Save(ctx context.Context, r *Record) (int64, error)
	Find(ctx context.Context, id int64) (*Record, error)
	ThisWeek(ctx context.Context, userId int64) ([]Record, error)
	ThisMonth(ctx context.Context, userId int64) ([]Record, error)
	ThisYear(ctx context.Context, userId int64) ([]Record, error)
	ThisMonthSum(ctx context.Context, userId int64) (float64, error)
	Count(ctx context.Context, userId int64) (int64, error)
}

var layouts = []string{"02-01-2006", "02:01:2006", "02 01 2006"}

func parseDate(s string) (time.Time, error) {
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("could not parse date")
}

func NewRecord(userId int64, msg string) (*Record, error) {
	reader := strings.NewReader(msg)
	record := &Record{UserId: userId}
	fmt.Fscanf(reader, "%f %s ", &record.Amount, &record.Category)

	if record.Amount <= 0 {
		return nil, errors.New("incorrect amount value")
	}
	if record.Category == "" {
		return nil, errors.New("empty category")
	}

	bytes, _ := io.ReadAll(reader)
	var err error
	record.Date, err = parseDate(string(bytes))
	if err != nil {
		return nil, err
	}

	return record, nil
}
