package repository

import (
	"context"
	"fmt"
	"time"

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
	r.logger.WithField("update_deposit", accountID).Debug("Prossesing deposit...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ct, err := r.db.Exec(ctx, queryDeposit, amount, accountID)
	if err != nil {
		r.logger.WithError(err).Error("Failed deposit")
		return fmt.Errorf("deposit failed: %w", err)
	}
	if ct.RowsAffected() == 0 {
		r.logger.WithError(err).Error("Failed not found account or locked")
		return fmt.Errorf("account %d not found or locked", accountID)
	}
	r.logger.Info("Successfully deposit")
	return nil
}

func (r *TransactionRepository) Withdraw(accountID int64, amount float64) error {
	r.logger.WithField("update_withdraw", accountID).Debug("Prossesing withdraw...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ct, err := r.db.Exec(ctx, queryWithdraw, amount, accountID)
	if err != nil {
		r.logger.WithError(err).Error("Failed withdraw")
		return fmt.Errorf("withdraw failed: %w", err)
	}
	if ct.RowsAffected() == 0 {
		r.logger.WithError(err).Error("Failed withdraw: insufficient funds or account locked")
		return fmt.Errorf("withdraw failed: insufficient funds or account locked")
	}

	r.logger.Info("Successfully withdraw")
	return nil
}

func (r *TransactionRepository) Transfer(fromAccountID, toAccountID int64, amount float64) error {
	r.logger.WithField("transfering", fromAccountID).Debug("Prossesing transfer...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.logger.WithError(err).Error("Failed begin tx failed")
		return fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback(ctx)

	ct, err := tx.Exec(ctx, queryWithdraw, amount, fromAccountID)
	if err != nil {
		r.logger.WithError(err).Error("Failed withdraw transfer")
		return fmt.Errorf("withdraw in transfer failed: %w", err)
	}
	if ct.RowsAffected() == 0 {
		r.logger.WithError(err).Error("Failed to transfer: insufficient funds or account locked")
		return fmt.Errorf("transfer failed: insufficient funds or account locked")
	}

	ct, err = tx.Exec(ctx, queryDeposit, amount, toAccountID)
	if err != nil {
		r.logger.WithError(err).Error("Failed to deposit in transfer")
		return fmt.Errorf("deposit in transfer failed: %w", err)
	}
	if ct.RowsAffected() == 0 {
		r.logger.WithError(err).Error("Failed to transfer: target account not found or locked")
		return fmt.Errorf("transfer failed: target account not found or locked")
	}

	if err := tx.Commit(ctx); err != nil {
		r.logger.WithError(err).Error("Failed to commit transfer")
		return fmt.Errorf("commit transfer failed: %w", err)
	}

	r.logger.Info("Successfully transfer")
	return nil
}

func (r *TransactionRepository) SaveTransaction(txn *entity.Transaction) error {
	r.logger.WithField("Saving transaction...", txn).Debug("Prossesing save transaction...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, querySave,
		txn.AccountID,
		txn.Amount,
		txn.TransactionType,
		txn.CreatedAt,
	).Scan(&txn.ID)

	if err != nil {
		r.logger.WithError(err).Error("Failed save transaction")
		return fmt.Errorf("save transaction failed: %w", err)
	}

	return nil
}

func (r *TransactionRepository) ListTransactions(accountID int64) ([]entity.Transaction, error) {
	r.logger.WithField("Listing transactions...", accountID).Debug("Prossesing list transactions...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, queryList, accountID)
	if err != nil {
		r.logger.WithError(err).Error("Failed list transactions")
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

	r.logger.Info("Successfully fetched list transactions")
	return transactions, nil
}
