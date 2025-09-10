package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/serikdev/CashFlow/internal/entity"
	"github.com/sirupsen/logrus"
)

type TransactionRepo interface {
	ListTransactions(accountID int64) ([]entity.Transaction, error)
}

type Producer interface {
	Publish(topic, key string, value []byte) error
}

type TransactionService struct {
	repo     TransactionRepo
	producer Producer
	logger   *logrus.Entry
}

func NewTransactionService(producer Producer, logger *logrus.Entry) *TransactionService {
	return &TransactionService{
		producer: producer,
		logger:   logger,
	}
}

func (s *TransactionService) Deposit(accountID int64, amount float64) (*entity.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("deposit amount must be greater than zero")
	}

	event := entity.TransactionEvent{
		AccountID:       accountID,
		Amount:          amount,
		TransactionType: "deposit",
		CreatedAt:       time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("error to marshal deposit event: %w", err)
	}

	if err := s.producer.Publish("account_deposit", fmt.Sprintf("%d", accountID), data); err != nil {
		return nil, fmt.Errorf("error to publish deposit event: %w", err)
	}
	return &entity.Transaction{
		AccountID:       int(accountID),
		Amount:          amount,
		TransactionType: "deposit",
		CreatedAt:       event.CreatedAt,
	}, nil
}

func (s *TransactionService) Withdraw(accountID int64, amount float64) (*entity.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("withdraw amount must be greater than zero")
	}

	event := entity.TransactionEvent{
		AccountID:       accountID,
		Amount:          amount,
		TransactionType: "withdraw",
		CreatedAt:       time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal withdraw event: %w", err)
	}

	if err := s.producer.Publish("account-withdraw", fmt.Sprintf("%d", accountID), data); err != nil {
		return nil, fmt.Errorf("failed to publish withdraw event: %w", err)
	}

	return &entity.Transaction{
		AccountID:       int(accountID),
		Amount:          amount,
		TransactionType: "withdraw",
		CreatedAt:       event.CreatedAt,
	}, nil
}

func (s *TransactionService) Transfer(fromAccountID, toAccountID int64, amount float64) (*entity.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("transfer amount must be greater than zero")
	}
	if fromAccountID == toAccountID {
		return nil, errors.New("cannot transfer to the same account")
	}

	event := entity.TransactionEvent{
		AccountID:       fromAccountID,
		RelatedAccount:  &toAccountID,
		Amount:          amount,
		TransactionType: "transfer",
		CreatedAt:       time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer event: %w", err)
	}

	if err := s.producer.Publish("account-transfer", fmt.Sprintf("%d", fromAccountID), data); err != nil {
		return nil, fmt.Errorf("failed to publish transfer event: %w", err)
	}

	return &entity.Transaction{
		AccountID:       int(fromAccountID),
		Amount:          amount,
		TransactionType: "transfer",
		CreatedAt:       event.CreatedAt,
	}, nil
}

func (s *TransactionService) ListTransactions(accountID int64) ([]entity.Transaction, error) {
	return s.repo.ListTransactions(accountID)
}
