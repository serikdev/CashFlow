package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/serikdev/CashFlow/internal/entity"
	"github.com/sirupsen/logrus"
)

type AccountRepo struct {
	db     *pgxpool.Pool
	logger *logrus.Entry
}

func NewAccountRepository(db *pgxpool.Pool, logger *logrus.Entry) *AccountRepo {
	return &AccountRepo{
		db:     db,
		logger: logger,
	}
}

const (
	createQuery = `
		INSERT INTO accounts(balance, currency, is_locked, created_at, deleted_at)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, balance, currency, is_locked, created_at, deleted_at
	`

	getByIDQuery = `
		SELECT id, balance, currency, is_locked, created_at, deleted_at
		FROM accounts
		WHERE id=$1
	`

	deleteQuery = `DELETE FROM accounts WHERE id=$1`

	listQuery = `
		SELECT id, balance, currency, is_locked, created_at, deleted_at
		FROM accounts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	countQuery = `SELECT COUNT(*) FROM accounts`
)

func (r *AccountRepo) Create(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	r.logger.WithField("account_balance", account.Balance).Debug("Creating account")

	var createAccount entity.Account
	err := r.db.QueryRow(
		ctx,
		createQuery,
		account.Balance,
		account.Currency,
		account.IsLocked,
		account.CreatedAt,
		account.DeletedAt,
	).Scan(
		&createAccount.ID,
		&createAccount.Balance,
		&createAccount.Currency,
		&createAccount.IsLocked,
		&createAccount.CreatedAt,
		&createAccount.DeletedAt,
	)

	if err != nil {
		r.logger.WithError(err).Error("Failed to create account in DB")
		return nil, fmt.Errorf("error to create account in DB: %w", err)
	}

	r.logger.WithField("account_id", createAccount.ID).Info("Successfully created account in DB")
	return &createAccount, nil
}

func (r *AccountRepo) GetByID(ctx context.Context, id int64) (*entity.Account, error) {
	r.logger.WithField("account_id", id).Debug("Fetching account ID")

	var account entity.Account

	err := r.db.QueryRow(ctx, getByIDQuery, id).Scan(
		&account.ID,
		&account.Balance,
		&account.Currency,
		&account.IsLocked,
		&account.CreatedAt,
		&account.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.WithField("account_id", account.ID).Warn("Account not found")
			return nil, fmt.Errorf("account with %d not found", id)
		}
		r.logger.WithFields(logrus.Fields{
			"account_id": id,
			"error":      err,
		}).Error("Failed to fetch account from DB")
		return nil, fmt.Errorf("error to fetch account: %w", err)
	}
	r.logger.WithField("account_id", account.ID).Info("Account fetched successfully from DB")
	return &account, nil
}

func (r *AccountRepo) Delete(ctx context.Context, id int64) error {
	r.logger.WithField("account_id", id).Debug("Removing account")

	cmdTag, err := r.db.Exec(ctx, deleteQuery, id)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"account_id": id,
			"error":      err,
		}).Error("Failed to delete account in DB")
		return fmt.Errorf("error  delete account in DB: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		r.logger.WithField("account_id", id).Warn("No rows were deleted")
		return fmt.Errorf("no rows account found with: %d", id)
	}

	r.logger.WithField("acount_id", id).Info("Account deleted successfully from database")
	return nil
}

func (r *AccountRepo) List(ctx context.Context, offset, limit int) ([]entity.Account, int, error) {
	r.logger.WithFields(logrus.Fields{
		"offset": offset,
		"limit":  limit,
	}).Debug("Fetching account list")

	var totalCount int
	err := r.db.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		r.logger.WithError(err).Error("Failed to count account")
		return nil, 0, fmt.Errorf("failed to count account: %w", err)
	}

	rows, err := r.db.Query(ctx, listQuery, limit, offset)
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"offset": offset,
			"limit":  limit,
			"error":  err,
		}).Error("Failed to fetch account from database")
		return nil, 0, fmt.Errorf("failed to fetch account from db: %w", err)
	}
	defer rows.Close()

	var accounts []entity.Account
	for rows.Next() {
		var account entity.Account
		err := rows.Scan(
			&account.ID,
			&account.Balance,
			&account.Currency,
			&account.IsLocked,
			&account.CreatedAt,
			&account.DeletedAt,
		)
		if err != nil {
			r.logger.WithError(err).Error("Failed to scan account row")
			return nil, 0, fmt.Errorf("failed to scan account row: %w", err)
		}
		accounts = append(accounts, account)

		if err = rows.Err(); err != nil {
			r.logger.WithError(err).Error("Error ocured during rows iteration")
			return nil, 0, fmt.Errorf("error ocured during rows iteration: %w", err)
		}

		r.logger.WithFields(logrus.Fields{
			"offset":      offset,
			"limit":       limit,
			"found_count": len(accounts),
			"total_count": totalCount,
		}).Info("News post list fetched Successfully from database")
	}
	return accounts, totalCount, nil
}
