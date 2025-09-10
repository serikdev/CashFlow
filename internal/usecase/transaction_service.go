package usecase

import (
	"encoding/json"
	"time"

	"github.com/serikdev/CashFlow/internal/entity"
	"github.com/serikdev/CashFlow/internal/kafka"
)

type TransactionUsecase struct {
	producer *kafka.Producer
}

func NewTransactionService(producer *kafka.Producer) *TransactionUsecase {
	return &TransactionUsecase{producer: producer}
}

func (u *TransactionUsecase) Deposit(accountID int64, amount float64) error {
	event := entity.TransactionEvent{
		Type:      "deposit",
		AccountID: accountID,
		Amount:    amount,
		Currency:  "TMT",
		CreatedAt: time.Now().UTC(),
	}
	msg, _ := json.Marshal(event)
	return u.producer.Publish("account-deposit", string(rune(accountID)), msg)
}

func (u *TransactionUsecase) Withdraw(accountID int64, amount float64) error {
	event := entity.TransactionEvent{
		Type:      "withdrawal",
		AccountID: accountID,
		Amount:    amount,
		Currency:  "TMT",
		CreatedAt: time.Now().UTC(),
	}
	msg, _ := json.Marshal(event)
	return u.producer.Publish("account-withdraw", string(rune(accountID)), msg)
}

func (u *TransactionUsecase) Transfer(fromID, toID int64, amount float64) error {
	event := entity.TransactionEvent{
		Type:        "transfer",
		AccountID:   fromID,
		ToAccountID: toID,
		Amount:      amount,
		Currency:    "TMT",
		CreatedAt:   time.Now().UTC(),
	}
	msg, _ := json.Marshal(event)
	return u.producer.Publish("account-transfer", string(rune(fromID)), msg)
}
