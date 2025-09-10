package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/serikdev/CashFlow/internal/entity"
	"github.com/sirupsen/logrus"
)

type TransactionRepository interface {
	ApplyTransaction(ctx context.Context, event entity.TransactionEvent) error
}

type transactionRepository struct {
	db     *pgxpool.Pool
	logger *logrus.Entry
}

func NewTransactionRepository(db *pgxpool.Pool, logger *logrus.Entry) TransactionRepository {
	return &transactionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *transactionRepository) ApplyTransaction(ctx context.Context, event entity.TransactionEvent) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	switch event.Type {
	case "deposit":
		_, err = tx.Exec(ctx,
			`UPDATE accounts SET balance = balance + $1 WHERE id = $2`,
			event.Amount, event.AccountID)
	case "withdrawal":
		_, err = tx.Exec(ctx,
			`UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1`,
			event.Amount, event.AccountID)
	case "transfer":
		_, err = tx.Exec(ctx,
			`UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1`,
			event.Amount, event.AccountID)
		if err == nil {
			_, err = tx.Exec(ctx,
				`UPDATE accounts SET balance = balance + $1 WHERE id = $2`,
				event.Amount, event.ToAccountID)
		}
	}

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO transactions (account_id, amount, transaction_type) VALUES ($1, $2, $3)`,
		event.AccountID, event.Amount, event.Type)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
