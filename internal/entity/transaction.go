package entity

import "time"

type Transaction struct {
	ID              int        `json:"id"`
	AccountID       int        `json:"account_id"`
	Amount          float64    `json:"amount"`
	TransactionType string     `json:"transaction_type"`
	CreatedAt       time.Time  `json:"created_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}
