package entity

import "time"

type TransactionEvent struct {
	Type        string    `json:"type"` // deposit, withdrawal, transfer
	AccountID   int64     `json:"account_id"`
	ToAccountID int64     `json:"to_account_id,omitempty"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
}
