package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const userId int64 = 123

func Test_OnCorrectMessage_ConstructorShouldParseString(t *testing.T) {
	msg := "12.21 food 12 04 2001"
	record, err := NewRecord(userId, msg)

	assert.NoError(t, err)
	assert.Equal(t, 12.21, record.Amount)
	assert.Equal(t, "food", record.Category)

	testDate, _ := time.Parse("02 01 2006", "12 04 2001")

	assert.Equal(t, testDate, record.Date)
}

func Test_OnInCorrectMessage_ConstructorShouldReturnError(t *testing.T) {
	msg := "something 12.21 bad"
	_, err := NewRecord(userId, msg)

	assert.Error(t, err)
}

func Test_ConstructorShouldParseDifferentDateFormats(t *testing.T) {
	testDate, _ := time.Parse("02 01 2006", "15 05 2001")

	record, err := NewRecord(userId, "12.21 food 15 05 2001")
	assert.NoError(t, err)
	assert.Equal(t, testDate, record.Date)

	record, err = NewRecord(userId, "12.21 food 15-05-2001")
	assert.NoError(t, err)
	assert.Equal(t, testDate, record.Date)

	record, err = NewRecord(userId, "12.21 food 15:05:2001")
	assert.NoError(t, err)
	assert.Equal(t, testDate, record.Date)
}
