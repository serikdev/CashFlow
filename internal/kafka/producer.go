package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Producer struct {
	writer *kafka.Writer
	logger *logrus.Entry
}

func NewProducer(brokers []string, logger *logrus.Entry) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
		logger: logger,
	}
}

func (p *Producer) Publish(topic, key string, message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: message,
	})

	if err != nil {
		p.logger.WithError(err).Error("failed to publish Kafka message")
		return err
	}

	p.logger.Info("Kafka message published to %s: %s", topic, string(message))
	return nil
}
