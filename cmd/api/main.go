package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/serikdev/CashFlow/docs"

	"github.com/serikdev/CashFlow/internal/adapter/repository"
	"github.com/serikdev/CashFlow/internal/config"
	"github.com/serikdev/CashFlow/internal/kafka"
	"github.com/serikdev/CashFlow/internal/port/rest"
	"github.com/serikdev/CashFlow/internal/port/rest/handler"
	"github.com/serikdev/CashFlow/internal/usecase"
	"github.com/serikdev/CashFlow/pkg/database"
	"github.com/serikdev/CashFlow/pkg/logger"
)

// @title CashFlow API
// @version 1.0
// @description API для управления счетами и транзакциями (Clean Architecture + Kafka).
// @host localhost:8080
// @BasePath /api
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.LoadConfig()
	log := logger.NewLogger()

	log.WithFields(map[string]interface{}{
		"db_host":   cfg.DBConfig.Host,
		"db_port":   cfg.DBConfig.Port,
		"db_name":   cfg.DBConfig.Name,
		"log_level": cfg.LoggerConfig.LogLevel,
	}).Info("Starting server with config")

	db, err := database.NewPool(ctx, cfg.DBConfig, cfg, log)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect database")
	}
	defer db.Close()

	producer := kafka.NewProducerImpl(cfg.KafkaConfig.Brokers, log)
	defer producer.Close()

	accountRepo := repository.NewAccountRepository(db, log)
	transactionRepo := repository.NewTransactionRepository(db, log)

	accountService := usecase.NewAccountService(accountRepo, log)

	transactionService := usecase.NewTransactionService(usecase.TransactionServiceDeps{
		TransactionRepo: transactionRepo,
		AccountRepo:     accountRepo,
		Producer:        producer,
		Logger:          log,
	})

	transactionService.SetRepo(transactionRepo)

	baseHandler := handler.NewBaseHandler(log)

	accountHandler := handler.NewAccountHandler(&baseHandler, accountService, log)
	transactionHandler := handler.NewTransactionHandler(&baseHandler, transactionService, log)

	handlers := rest.Handlers{
		AccountHandler:     accountHandler,
		TransactionHandler: transactionHandler,
	}

	router := rest.NewRouter(&handlers)

	// Start Kafka Consumers
	go func() {
		depositConsumer := kafka.NewConsumerImpl(cfg.KafkaConfig.Brokers, "account-deposit", "cashflow-group", transactionRepo, log)
		if err := depositConsumer.Run(ctx); err != nil {
			log.WithError(err).Fatal("Deposit consumer failed")
		}
	}()
	go func() {
		withdrawConsumer := kafka.NewConsumerImpl(cfg.KafkaConfig.Brokers, "account-withdraw", "cashflow-group", transactionRepo, log)
		if err := withdrawConsumer.Run(ctx); err != nil {
			log.WithError(err).Fatal("Withdraw consumer failed")
		}
	}()
	go func() {
		transferConsumer := kafka.NewConsumerImpl(cfg.KafkaConfig.Brokers, "account-transfer", "cashflow-group", transactionRepo, log)
		if err := transferConsumer.Run(ctx); err != nil {
			log.WithError(err).Fatal("Transfer consumer failed")
		}
	}()

	// HTTP Server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Server failed to start")
		}
	}()

	<-stop
	log.Info("Shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.WithError(err).Error("Server forced to shutdown")
	} else {
		log.Info("Server exited gracefully")
	}

	log.Info("Server stopped")
}
