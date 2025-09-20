package usecase

import (
	"context"
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
type AccountRepository interface {
	GetByID(ctx context.Context, accountID int64) (*entity.Account, error)
}

type Producer interface {
	Publish(topic, key string, value []byte) error
}

type TransactionServiceDeps struct {
	TransactionRepo TransactionRepo
	AccountRepo     AccountRepo
	Producer        Producer
	Logger          *logrus.Entry
}

type TransactionService struct {
	transacRepo TransactionRepo
	accountRepo AccountRepo
	producer    Producer
	logger      *logrus.Entry
}

func NewTransactionService(deps TransactionServiceDeps) *TransactionService {
	return &TransactionService{
		transacRepo: deps.TransactionRepo,
		accountRepo: deps.AccountRepo,
		producer:    deps.Producer,
		logger:      deps.Logger,
	}
}

func (s *TransactionService) Deposit(ctx context.Context, accountID int64, amount float64) (*entity.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("deposit amount must be greater than zero")
	}
	_, err := s.checkAccountActive(ctx, accountID)
	if err != nil {
		return nil, err
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

	if err := s.producer.Publish("account-deposit", fmt.Sprintf("%d", accountID), data); err != nil {
		return nil, fmt.Errorf("error to publish deposit event: %w", err)
	}
	return &entity.Transaction{
		AccountID:       int(accountID),
		Amount:          amount,
		TransactionType: "deposit",
		CreatedAt:       event.CreatedAt,
	}, nil
}

func (s *TransactionService) Withdraw(ctx context.Context, accountID int64, amount float64) (*entity.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("withdraw amount must be greater than zero")
	}

	_, err := s.checkAccountActive(ctx, accountID)
	if err != nil {
		return nil, err
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

func (s *TransactionService) Transfer(ctx context.Context, fromAccountID, toAccountID int64, amount float64) (*entity.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("transfer amount must be greater than zero")
	}
	if fromAccountID == toAccountID {
		return nil, errors.New("cannot transfer to the same account")
	}

	_, err := s.checkAccountActive(ctx, fromAccountID)
	if err != nil {
		return nil, err
	}
	_, err = s.checkAccountActive(ctx, toAccountID)
	if err != nil {
		return nil, err
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
	return s.transacRepo.ListTransactions(accountID)
}

func (s *TransactionService) SetRepo(TransacRepo TransactionRepo) {
	s.transacRepo = TransacRepo
}

func (s *TransactionService) checkAccountActive(ctx context.Context, accountID int64) (*entity.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	if account.IsLocked {
		return nil, errors.New("account is locked")
	}

	if account.DeletedAt != nil {
		return nil, errors.New("account is deleted")
	}
	return account, nil

}
