package entity

import "time"

type TransactionEvent struct {
	AccountID       int64     `json:"account_id"`
	RelatedAccount  *int64    `json:"related_account,omitempty"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	CreatedAt       time.Time `json:"created_at"`
}
