package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/serikdev/CashFlow/internal/entity"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type TransactionRepository interface {
	Deposit(accountID int64, amount float64) error
	Withdraw(accountID int64, amount float64) error
	Transfer(fromAccountID, toAccountID int64, amount float64) error
	SaveTransaction(tx *entity.Transaction) error
}

type TransactionEvent struct {
	AccountID       int64     `json:"account_id"`
	RelatedAccount  *int64    `json:"related_account,omitempty"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	CreatedAt       time.Time `json:"created_at"`
}

type ConsumerImpl struct {
	reader     *kafka.Reader
	logger     *logrus.Entry
	repository TransactionRepository
}

func NewConsumerImpl(brokers []string, topic, groupID string, repo TransactionRepository, logger *logrus.Entry) *ConsumerImpl {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		GroupTopics: []string{"account-deposit", "account-withdraw", "account-transfer"},
		Topic:       topic,
		MinBytes:    10e3,
		MaxBytes:    10e6,
	})

	return &ConsumerImpl{
		reader:     r,
		logger:     logger.WithField("topic", topic),
		repository: repo,
	}
}

func (c *ConsumerImpl) Run(ctx context.Context) error {
	c.logger.Info("Consumer started")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {

			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				c.logger.Info("Consumer context cancelled, shutting down gracefully")
				return nil
			}

			c.logger.WithError(err).Error("Failed to read message, will retry...")
			continue
		}

		var event TransactionEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			c.logger.WithError(err).Error("Failed to unmarshal event")
			continue
		}

		c.logger.WithFields(logrus.Fields{
			"key":   string(m.Key),
			"value": string(m.Value),
		}).Info("Message received")

		switch event.TransactionType {
		case "deposit":
			err = c.repository.Deposit(event.AccountID, event.Amount)
		case "withdraw":
			err = c.repository.Withdraw(event.AccountID, event.Amount)
		case "transfer":
			if event.RelatedAccount == nil {
				err = fmt.Errorf("related account is nil for transfer")
			} else {
				err = c.repository.Transfer(event.AccountID, *event.RelatedAccount, event.Amount)
			}
		default:
			c.logger.Warnf("unknown transaction type: %s", event.TransactionType)
			continue
		}

		if err != nil {
			c.logger.WithError(err).Error("Failed to process transaction")
			continue
		}

		tx := &entity.Transaction{
			AccountID:       int(event.AccountID),
			Amount:          event.Amount,
			TransactionType: event.TransactionType,
			CreatedAt:       event.CreatedAt,
		}
		if err := c.repository.SaveTransaction(tx); err != nil {
			c.logger.WithError(err).Error("Failed to save transaction")
		}
	}
}

func (c *ConsumerImpl) Close() error {
	return c.reader.Close()
}
