package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/serikdev/CashFlow/internal/adapter/repository"
	"github.com/serikdev/CashFlow/internal/entity"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	reader     *kafka.Reader
	logger     *logrus.Entry
	repository repository.TransactionRepository
}

func NewConsumer(brokers []string, topic, groupID string, logger *logrus.Entry, repo repository.TransactionRepository) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		logger:     logger,
		repository: repo,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.logger.WithError(err).Error("Kafka read error")
			continue
		}

		var event entity.TransactionEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			c.logger.WithError(err).Error("Invalid event format")
			continue
		}

		c.logger.Infof("Received Kafka event: %s", string(m.Value))

		if err := c.repository.ApplyTransaction(ctx, event); err != nil {
			c.logger.WithError(err).Error("Failed to apply transaction")
		}
	}
}
