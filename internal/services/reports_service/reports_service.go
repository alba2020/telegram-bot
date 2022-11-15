package reportsservice

import (
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
)

type Config interface {
	BrokersList() []string
	Topic() string
}

type ReportsService struct {
	producer    sarama.SyncProducer
	topic       string
	brokersList []string
}

func New(cfg Config) (*ReportsService, error) {
	logger.Infof("Kafka brokers: %s", strings.Join(cfg.BrokersList(), ", "))

	producer, err := newSyncProducer(cfg.BrokersList())
	if err != nil {
		return nil, err
	}

	return &ReportsService{
		producer:    producer,
		topic:       cfg.Topic(),
		brokersList: cfg.BrokersList(),
	}, nil
}

func (svc *ReportsService) SendMessage(key, value string) (int32, int64, error) {
	msg := sarama.ProducerMessage{
		Topic: svc.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	return svc.producer.SendMessage(&msg)
}

func (svc *ReportsService) RequestReport(key, value string) error {
	p, o, err := svc.SendMessage(key, value)
	if err != nil {
		return err
	} else {
		logger.Debugf("Write message %s %s\n", key, value)
		logger.Debugf("Success: topic %s, offset:%d, partition: %d\n", svc.topic, o, p)
		return nil
	}
}

func (svc *ReportsService) Close() error {
	return svc.producer.Close()
}

func newSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	// Waits for all in-sync replicas to commit before responding.
	config.Producer.RequiredAcks = sarama.WaitForAll
	// The total number of times to retry sending a message (default 3).
	config.Producer.Retry.Max = 3
	// How long to wait for the cluster to settle between retries (default 100ms).
	config.Producer.Retry.Backoff = time.Millisecond * 250
	// idempotent producer has a unique producer ID and uses sequence IDs for each message,
	// allowing the broker to ensure, on a per-partition basis, that it is committing ordered messages with no duplication.
	//config.Producer.Idempotent = true
	if config.Producer.Idempotent {
		config.Producer.Retry.Max = 1
		config.Net.MaxOpenRequests = 1
	}
	//  Successfully delivered messages will be returned on the Successes channe
	config.Producer.Return.Successes = true
	// Generates partitioners for choosing the partition to send messages to (defaults to hashing the message key)
	_ = config.Producer.Partitioner

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, fmt.Errorf("starting Sarama producer: %w", err)
	}

	return producer, nil
}
