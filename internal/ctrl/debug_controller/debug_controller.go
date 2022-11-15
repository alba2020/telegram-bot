package debug_controller

import (
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
)

type DebugController struct {
	recordRepo domain.RecordRepository
}

func New(repo domain.RecordRepository) (*DebugController, error) {
	return &DebugController{
		recordRepo: repo,
	}, nil
}

func (c *DebugController) SaveCommand(msg *domain.Message) *domain.Message {
	record := domain.Record{
		UserId:   1223,
		Amount:   101,
		Category: "milk",
		Date:     time.Now(),
	}

	id, err := c.recordRepo.Save(msg.Context, &record)
	if err != nil {
		return msg.Reply(err.Error())
	} else {
		return msg.Reply(fmt.Sprintf("saved with id %d", id))
	}
}

func (c *DebugController) SumCommand(msg *domain.Message) *domain.Message {
	sum, _ := c.recordRepo.ThisMonthSum(msg.Context, msg.UserID)
	return msg.Reply(fmt.Sprintf("this month sum %f", sum))
}
