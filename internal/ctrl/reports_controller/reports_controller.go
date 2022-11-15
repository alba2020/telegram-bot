package reports_controller

import (
	"fmt"

	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
	reportsservice "gitlab.ozon.dev/albatros2002/telegram-bot/internal/services/reports_service"
)

type ReportsController struct {
	reportsService *reportsservice.ReportsService
}

func New(svc *reportsservice.ReportsService) *ReportsController {
	return &ReportsController{
		reportsService: svc,
	}
}

func (c *ReportsController) ReportCommand(msg *domain.Message) *domain.Message {
	k := fmt.Sprintf("%d", msg.UserID)
	v := "report"
	err := c.reportsService.RequestReport(k, v)
	if err != nil {
		logger.Fatal("Failed to send message:", err)
	}

	return msg.Reply("Generating report...")
}
