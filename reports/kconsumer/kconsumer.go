package kconsumer

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/albatros2002/telegram-bot/internal/logger"
)

var Assignor = "range"

type Config interface {
	Topic() string
	ConsumerGroup() string
	BrokersList() []string
}

func StartConsumerGroup(
	cfg Config,
	ctx context.Context,
	f func(*sarama.ConsumerMessage)) error {
	logger.Infof("Kafka brokers: %s", strings.Join(cfg.BrokersList(), ", "))

	consumerGroupHandler := NewConsumer(f)

	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	switch Assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
	case "round-robin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", Assignor)
	}

	// Create consumer group
	consumerGroup, err := sarama.NewConsumerGroup(
		cfg.BrokersList(), cfg.ConsumerGroup(), config)
	if err != nil {
		return fmt.Errorf("starting consumer group: %w", err)
	}

	err = consumerGroup.Consume(ctx, []string{cfg.Topic()}, &consumerGroupHandler)
	if err != nil {
		return fmt.Errorf("consuming via handler: %w", err)
	}
	return nil
}

// Consumer represents a Sarama consumer group consumer.
type Consumer struct {
	handler func(*sarama.ConsumerMessage)
}

func NewConsumer(f func(*sarama.ConsumerMessage)) Consumer {
	return Consumer{
		handler: f,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	logger.Debug("consumer - setup")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	logger.Debug("consumer - cleanup")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// printMessage(message)
		consumer.handler(message)
		session.MarkMessage(message, "")
	}

	return nil
}
