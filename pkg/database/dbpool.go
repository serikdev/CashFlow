package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/serikdev/CashFlow/internal/config"
	"github.com/sirupsen/logrus"
)

func NewPool(ctx context.Context, cfg config.DBConfig, logger *logrus.Entry) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	logger.WithFields(logrus.Fields{
		"host":    cfg.Host,
		"port":    cfg.Port,
		"user":    cfg.User,
		"name":    cfg.Name,
		"sslmode": cfg.SllMode,
	}).Info("Connecting Database")

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.WithError(err).Error("Failed to parse pool config")
		return nil, fmt.Errorf("error to parse pool config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.WithError(err).Error("Failed to pool DB")
		return nil, fmt.Errorf("error to pool db: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		logger.WithError(err).Error("Failed to ping DB")
		return nil, fmt.Errorf("error to ping DB: %w", err)
	}
	logger.Info("Successfully connected with DB")

	return pool, nil

}
