package recordmemoryrepo

import (
	"context"
	"errors"
	"sync"
	"time"

	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
)

type Model struct {
	mu      *sync.RWMutex
	lastId  int64
	records []domain.Record
}

func New() *Model {
	return &Model{
		mu:      &sync.RWMutex{},
		lastId:  0,
		records: make([]domain.Record, 0),
	}
}

func (m *Model) Save(ctx context.Context, rec *domain.Record) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastId++
	rec.Id = m.lastId
	m.records = append(m.records, *rec)
	return rec.Id, nil
}

func (m *Model) Find(ctx context.Context, id int64) (*domain.Record, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, record := range m.records {
		if record.Id == id {
			return &record, nil
		}
	}
	return nil, errors.New("not found")
}

// week month year

type Predicate func(record domain.Record) bool

func (m *Model) Filter(f Predicate) []domain.Record {
	m.mu.RLock()
	defer m.mu.RUnlock()

	found := []domain.Record{}
	for _, rec := range m.records {
		if f(rec) {
			found = append(found, rec)
		}
	}
	return found
}

func (m *Model) ThisWeek(ctx context.Context, userId int64) ([]domain.Record, error) {
	currentYear, currentWeek := time.Now().ISOWeek()
	found := m.Filter(func(rec domain.Record) bool {
		recYear, recWeek := rec.Date.ISOWeek()
		return rec.UserId == userId &&
			recYear == currentYear &&
			recWeek == currentWeek
	})
	return found, nil
}

func (m *Model) ThisMonth(ctx context.Context, userId int64) ([]domain.Record, error) {
	currentMonth := time.Now().Month()
	currentYear := time.Now().Year()
	found := m.Filter(func(rec domain.Record) bool {
		return rec.UserId == userId &&
			rec.Date.Month() == currentMonth &&
			rec.Date.Year() == currentYear
	})
	return found, nil
}

func (m *Model) ThisYear(ctx context.Context, userId int64) ([]domain.Record, error) {
	currentYear := time.Now().Year()
	found := m.Filter(func(rec domain.Record) bool {
		return rec.UserId == userId &&
			currentYear == rec.Date.Year()
	})
	return found, nil
}

func (m *Model) ThisMonthSum(ctx context.Context, userId int64) (float64, error) {
	currentMonth := time.Now().Month()
	currentYear := time.Now().Year()

	var total float64 = 0
	for _, rec := range m.records {
		if rec.UserId == userId &&
			rec.Date.Month() == currentMonth &&
			rec.Date.Year() == currentYear {
			total += rec.Amount
		}
	}

	return total, nil
}

func (m *Model) Count(ctx context.Context, userId int64) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	found := m.Filter(func(rec domain.Record) bool {
		return rec.UserId == userId
	})

	return int64(len(found)), nil
}
