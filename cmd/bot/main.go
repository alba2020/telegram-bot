package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/cache"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/config"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/ctrl/debug_controller"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/ctrl/info_controller"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/ctrl/records_controller"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/ctrl/reports_controller"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/model/commands"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/model/currency"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/model/messages"
	recordsservice "gitlab.ozon.dev/albatros2002/telegram-bot/internal/services/records_service"
	reportsserver "gitlab.ozon.dev/albatros2002/telegram-bot/internal/services/reports_server"
	reportsservice "gitlab.ozon.dev/albatros2002/telegram-bot/internal/services/reports_service"

	// recordmemoryrepo "gitlab.ozon.dev/albatros2002/telegram-bot/internal/repo/record_memory_repo"
	recordmemoryrepo "gitlab.ozon.dev/albatros2002/telegram-bot/internal/repo/record_memory_repo"
	recordpgrepo "gitlab.ozon.dev/albatros2002/telegram-bot/internal/repo/record_pg_repo"
	sessionpgrepo "gitlab.ozon.dev/albatros2002/telegram-bot/internal/repo/session_pg_repo"
)

var (
	port = 8080
	// develMode   = flag.Bool("devel", false, "development mode")
	// serviceName = flag.String("service", "fibonacci", "the name of our service")
)

func init() {
	logger.Init(logger.Development)
	logger.Info("Starting bot")

	initTracing()

	go func() {
		http.Handle("/metrics", promhttp.Handler())

		logger.Info("starting http server", port)
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			logger.Fatal("error starting http server", err)
		}
	}()
}

func initTracing() {
	cfg := jaegerconfig.Configuration{
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	const serviceName = "cute-telegram-bot"
	_, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("Cannot init tracing", err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		signal.Stop(exit)
		cancel()
	}()

	config, err := config.New("data/config.yaml")
	if err != nil {
		logger.Fatal("Config init failed:", err)
	}

	tgClient, err := tg.New(config, ctx)
	if err != nil {
		logger.Fatal("Telegram client init failed:", err)
	}

	var currencyService commands.CurrencyService

	recordMemRepo := recordmemoryrepo.New()
	// sessionMemRepo, _ := sessionmemoryrepo.New()
	_ = recordMemRepo

	recordPGRepo, err := recordpgrepo.New(config.PGDSN())
	if err != nil {
		logger.Fatal("Could not create pg record repo", err.Error())
	}

	sessionPGRepo, err := sessionpgrepo.New(config.PGDSN())
	if err != nil {
		logger.Fatal("Could not create pg session repo", err.Error())
	}

	// async service initialization
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		currencyService, _ = currency.New()
	}()
	wg.Wait()

	recordsService, _ := recordsservice.New(sessionPGRepo, recordPGRepo, currencyService)

	var infoController = info_controller.New()
	var recordsController = records_controller.New(
		recordPGRepo, sessionPGRepo, currencyService, recordsService)
	var debugController, _ = debug_controller.New(recordPGRepo)

	var reportsService, _ = reportsservice.New(config)
	var reportsController = reports_controller.New(reportsService)

	_ = recordsController
	_ = debugController

	cache, err := cache.New(128)
	if err != nil {
		panic("Cache not started")
	}

	cached := cache.WithCache
	invalidate := cache.Invalidate

	var router = domain.NewRouter(map[string]domain.CommandHandler{
		"/start": MW(infoController.StartCommand),
		"/ping":  MW(infoController.PingCommand),
		"/add":   MW(invalidate(recordsController.AddCommand)),
		"/week":  MW(cached(recordsController.WeekCommand)),
		"/month": MW(cached(recordsController.MonthCommand)),
		"/year":  MW(cached(recordsController.YearCommand)),

		"/select_currency": MW(recordsController.SelectCurrencyCommand),
		"/set_currency":    MW(recordsController.SetCurrencyCommand),

		"/limit": MW(recordsController.LimitCommand),

		"/save": MW(debugController.SaveCommand),
		"/sum":  MW(debugController.SumCommand),

		"/report": reportsController.ReportCommand,
	})

	// start gRPC server
	go reportsserver.StartReportsServer(tgClient)

	msgModel := messages.New(tgClient, router)
	tgClient.ListenUpdates(msgModel)

	logger.Debug("Main done")
}
