package recordmemoryrepo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
)

var testRecord domain.Record

func init() {
	testDate, _ := time.Parse("02 01 2006", "15 05 2001")
	testRecord = domain.Record{
		Id:       0,
		UserId:   2,
		Amount:   3.14,
		Category: "food",
		Date:     testDate,
	}
}

func Test_AfterAddingRecord_ShouldIncreaseSliceLength(t *testing.T) {
	storage := New()
	ctx := context.Background()
	assert.Equal(t, 0, len(storage.records))
	_, _ = storage.Save(ctx, &testRecord)
	assert.Equal(t, 1, len(storage.records))
}

func Test_AddedRecord_CanBeFoundById(t *testing.T) {
	ctx := context.Background()
	storage := New()
	_, _ = storage.Save(ctx, &testRecord)
	found, err := storage.Find(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, testRecord, *found)
}

func seedRecords(storage domain.RecordRepository) {
	ctx := context.Background()
	_, _ = storage.Save(ctx, &domain.Record{
		Id:       0,
		UserId:   2,
		Amount:   3.14,
		Category: "food",
		Date:     time.Now(),
	})
	_, _ = storage.Save(ctx, &domain.Record{
		Id:       0,
		UserId:   2,
		Amount:   4.14,
		Category: "food",
		Date:     time.Now(),
	})
	_, _ = storage.Save(ctx, &domain.Record{
		Id:       0,
		UserId:   2,
		Amount:   5.14,
		Category: "food",
		Date:     time.Now(),
	})
}

func Test_SearchEntriesForThisWeek(t *testing.T) {
	ctx := context.Background()
	storage := New()
	seedRecords(storage)
	thisWeek, _ := storage.ThisWeek(ctx, 2)
	assert.Equal(t, 3, len(thisWeek))
}

func Test_SearchEntriesByAmount(t *testing.T) {
	storage := New()
	seedRecords(storage)
	found := storage.Filter(func(rec domain.Record) bool {
		return rec.Amount > 4
	})
	assert.Equal(t, 2, len(found))
}
