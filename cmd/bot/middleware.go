package main

import (
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
)

func WithLogging(cmd domain.CommandHandler) domain.CommandHandler {
	return func(msg *domain.Message) *domain.Message {
		logger.Info("=== before command ===")
		response := cmd(msg)
		logger.Info(response.Text)
		logger.Info("=== after command ===")
		return response
	}
}

func WithTracing(cmd domain.CommandHandler) domain.CommandHandler {
	return func(msg *domain.Message) *domain.Message {
		span, ctx := opentracing.StartSpanFromContext(
			msg.Context,
			"executing command "+msg.Command,
		)
		defer span.Finish()
		msg.Context = ctx

		response := cmd(msg)

		return response
	}
}

var InFlightRequests = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "telegram_bot",
	Subsystem: "requests",
	Name:      "in_flight_requests_total",
})

var SummaryResponseTime = promauto.NewSummary(prometheus.SummaryOpts{
	Namespace: "telegram_bot",
	Subsystem: "requests",
	Name:      "summary_response_time_seconds",
	Objectives: map[float64]float64{
		0.5:  0.1,
		0.9:  0.01,
		0.99: 0.001,
	},
})

func WithMetrics(cmd domain.CommandHandler) domain.CommandHandler {
	return func(msg *domain.Message) *domain.Message {
		InFlightRequests.Inc()
		defer InFlightRequests.Dec()

		startTime := time.Now()

		response := cmd(msg)

		duration := time.Since(startTime)
		SummaryResponseTime.Observe(duration.Seconds())

		return response
	}
}

func MW(cmd domain.CommandHandler) domain.CommandHandler {
	return WithLogging(
		WithTracing(
			WithMetrics(
				cmd,
			),
		),
	)
}
