package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/serikdev/CashFlow/internal/entity"
	"github.com/sirupsen/logrus"
)

type AccountRepo interface {
	Create(ctx context.Context, account *entity.Account) (*entity.Account, error)
	GetByID(ctx context.Context, id int64) (*entity.Account, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, offset, limit int) ([]entity.Account, int, error)
}

type AccountService struct {
	repo   AccountRepo
	logger *logrus.Entry
}

func NewAccountService(repo AccountRepo, logger *logrus.Entry) *AccountService {
	return &AccountService{
		repo:   repo,
		logger: logger,
	}
}

func (s *AccountService) Create(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	now := time.Now()
	account.CreatedAt = now

	createAccount, err := s.repo.Create(ctx, account)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create account")
		return nil, fmt.Errorf("error creating account: %w", err)
	}

	s.logger.WithField("account_id", createAccount.ID).Info("Successfully created account")
	return createAccount, nil
}

func (s *AccountService) GetByID(ctx context.Context, id int64) (*entity.Account, error) {
	if id <= 0 {
		return nil, errors.New("Invalid account ID")
	}

	existingAccount, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"account_id": id,
			"error":      err,
		}).Error("Failed to fetch existing account ID")
		return nil, fmt.Errorf("error fetching existing account ID: %w", err)
	}
	return existingAccount, nil
}

func (s *AccountService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("Invalid ID news post")
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"account_id": id,
			"error":      err,
		}).Error("Failed to fetch account ID")
		return fmt.Errorf("failed to fetch account ID: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.WithFields(logrus.Fields{
			"account_id": id,
			"error":      err,
		}).Error("Failed to remove account")

		return fmt.Errorf("failed to remove account: %w", err)
	}

	s.logger.WithField("account_id", id).Info("Successfilly deleted")
	return nil
}

func (s *AccountService) List(ctx context.Context, page, limit int) ([]entity.Account, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	account, total, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"page":   page,
			"limit":  limit,
			"offset": offset,
			"error":  err,
		}).Error("Failed to fetch account list")
		return nil, 0, fmt.Errorf("error to fetch account: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"page":        page,
		"limit":       limit,
		"total_found": len(account),
		"total_count": total,
	}).Debug("Account list fetched Successfully")
	return account, total, nil
}
