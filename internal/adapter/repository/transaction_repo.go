package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/serikdev/CashFlow/internal/entity"
	"github.com/sirupsen/logrus"
)

type TransactionRepository struct {
	db     *pgxpool.Pool
	logger *logrus.Entry
}

func NewTransactionRepository(db *pgxpool.Pool, logger *logrus.Entry) *TransactionRepository {
	return &TransactionRepository{
		db:     db,
		logger: logger,
	}
}

const (
	query = `UPDATE accounts SET balance = balance + $1 WHERE id = $2 AND deleted_at IS NULL AND is_locked = FALSE`

	queryB = `
		UPDATE accounts 
		SET balance = balance - $1
		WHERE id = $2 AND deleted_at IS NULL AND is_locked = FALSE AND balance >= $1
	`
	queryWithdraw = `
		UPDATE accounts 
		SET balance = balance - $1
		WHERE id = $2 AND deleted_at IS NULL AND is_locked = FALSE AND balance >= $1
	`
	queryDeposit = `
		UPDATE accounts 
		SET balance = balance + $1
		WHERE id = $2 AND deleted_at IS NULL AND is_locked = FALSE
	`
	querySave = `
		INSERT INTO transactions (account_id, amount, transaction_type, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	queryList = `
		SELECT id, account_id, amount, transaction_type, created_at, deleted_at
		FROM transactions
		WHERE account_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 100
	`
)

func (r *TransactionRepository) Deposit(accountID int64, amount float64) error {
	ctx := context.Background()

	ct, err := r.db.Exec(ctx, query, amount, accountID)
	if err != nil {
		return fmt.Errorf("deposit failed: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("account %d not found or locked", accountID)
	}
	return nil
}

func (r *TransactionRepository) Withdraw(accountID int64, amount float64) error {
	ctx := context.Background()

	ct, err := r.db.Exec(ctx, queryB, amount, accountID)
	if err != nil {
		return fmt.Errorf("withdraw failed: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("withdraw failed: insufficient funds or account locked")
	}
	return nil
}

func (r *TransactionRepository) Transfer(fromAccountID, toAccountID int64, amount float64) error {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback(ctx)

	ct, err := tx.Exec(ctx, queryWithdraw, amount, fromAccountID)
	if err != nil {
		return fmt.Errorf("withdraw in transfer failed: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("transfer failed: insufficient funds or account locked")
	}

	ct, err = tx.Exec(ctx, queryDeposit, amount, toAccountID)
	if err != nil {
		return fmt.Errorf("deposit in transfer failed: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("transfer failed: target account not found or locked")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transfer failed: %w", err)
	}

	return nil
}

func (r *TransactionRepository) SaveTransaction(txn *entity.Transaction) error {
	ctx := context.Background()

	err := r.db.QueryRow(ctx, querySave,
		txn.AccountID,
		txn.Amount,
		txn.TransactionType,
		txn.CreatedAt,
	).Scan(&txn.ID)

	if err != nil {
		return fmt.Errorf("save transaction failed: %w", err)
	}

	return nil
}

func (r *TransactionRepository) ListTransactions(accountID int64) ([]entity.Transaction, error) {
	ctx := context.Background()

	rows, err := r.db.Query(ctx, queryList, accountID)
	if err != nil {
		return nil, fmt.Errorf("list transactions failed: %w", err)
	}
	defer rows.Close()

	var transactions []entity.Transaction
	for rows.Next() {
		var t entity.Transaction
		if err := rows.Scan(
			&t.ID,
			&t.AccountID,
			&t.Amount,
			&t.TransactionType,
			&t.CreatedAt,
			&t.DeletedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
