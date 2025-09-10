package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type ProducerImpl struct {
	writer *kafka.Writer
	logger *logrus.Entry
}

func NewProducer(brokers []string, logger *logrus.Entry) *ProducerImpl {
	return &ProducerImpl{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
		},
		logger: logger,
	}
}

func (p *ProducerImpl) Publish(topic string, key string, value []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		p.logger.WithError(err).Errorf("failed to publish message to topic=%s", topic)
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"topic": topic,
		"key":   key,
	}).Info("Message published to Kafka")

	return nil
}

func (p *ProducerImpl) Close() error {
	return p.writer.Close()
}
