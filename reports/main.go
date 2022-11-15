package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/config"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/domain"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
	recordpgrepo "gitlab.ozon.dev/albatros2002/telegram-bot/internal/repo/record_pg_repo"
	"gitlab.ozon.dev/albatros2002/telegram-bot/reports/kconsumer"
	pb "gitlab.ozon.dev/albatros2002/telegram-bot/reports/reports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config interface {
	PGDSN() string
	ReportsServerAddress() string
	BrokersList() []string
	ConsumerGroup() string
	Topic() string
}

var cfg Config
var recordPGRepo domain.RecordRepository

func init() {
	var err error
	cfg, err = config.New("../data/config.yaml")
	if err != nil {
		logger.Fatal("Config init failed:", err)
	}

	recordPGRepo, err = recordpgrepo.New(cfg.PGDSN())
	if err != nil {
		logger.Fatal("Could not create pg record repo", err.Error())
	}
}

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := kconsumer.StartConsumerGroup(cfg, ctx, handleMessage); err != nil {
		logger.Fatal(err)
	}

	<-ctx.Done()
}

func handleMessage(msg *sarama.ConsumerMessage) {
	logger.Debugf("New message received from topic:%s, offset:%d, partition:%d, key:%s,"+
		" value:%s\n",
		msg.Topic, msg.Offset, msg.Partition, string(msg.Key), string(msg.Value))

	// Emulate Work loads
	// time.Sleep(1 * time.Second)
	logger.Debug("Successful to read message: ", string(msg.Value))

	userId, _ := strconv.Atoi(string(msg.Key))
	count, _ := recordPGRepo.Count(context.Background(), int64(userId))
	report := fmt.Sprintf("total records in DB: %d", count)

	dial(int64(userId), report)
}

func dial(userId int64, msg string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(cfg.ReportsServerAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Did not connect", err)
	}
	defer conn.Close()
	client := pb.NewReporterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.SendReport(
		ctx,
		&pb.ReportRequest{
			UserId: userId, Message: msg,
		},
	)

	if err != nil {
		logger.Fatal("Could not send report", err)
	}
	logger.Debugf("Send report status: %s", response.GetStatus())
}
