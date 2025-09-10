package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/serikdev/CashFlow/internal/adapter/repository"
	"github.com/serikdev/CashFlow/internal/config"
	"github.com/serikdev/CashFlow/internal/port/rest"
	"github.com/serikdev/CashFlow/internal/port/rest/handler"
	"github.com/serikdev/CashFlow/internal/usecase"
	"github.com/serikdev/CashFlow/pkg/database"
	"github.com/serikdev/CashFlow/pkg/logger"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	config := config.LoadConfig()
	logger := logger.NewLogger()

	// Adding the log at startup
	logger.WithFields(map[string]interface{}{
		"db_host":   config.DBConfig.Host,
		"db_port":   config.DBConfig.Port,
		"db_name":   config.DBConfig.Name,
		"log_level": config.LoggerConfig.LogLevel,
	}).Info("Starting server with safe config")

	db, err := database.NewPool(ctx, config.DBConfig, *logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect database")
	}
	defer db.Close()

	// Initializing layers
	accountRepo := repository.NewAccountRepository(db, logger)
	transactionRepo := repository.NewTransactionRepository(db, logger)

	// ==== Usecases ====
	accountService := usecase.NewAccountService(accountRepo, logger)
	transactionService := usecase.NewTransactionService(transactionRepo, accountRepo, logger)

	// ==== Handlers ====
	baseHandler := handler.NewBaseHandler(logger)

	accountHandler := handler.NewAccountHandler(&baseHandler, accountService, logger)
	transactionHandler := handler.NewTransactionHandler(&baseHandler, transactionService, logger)

	handlers := rest.Handlers{
		AccountHandler:     accountHandler,
		TransactionHandler: transactionHandler,
	}

	router := rest.NewRouter(&handlers)
	// HTTP server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Gracefull shutdowns
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Server failed to start")
		}
	}()

	<-stop
	logger.Info("SHutting down server.....")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := server.Shutdown(ctxShutdown); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	} else {
		logger.Info("Server exited gracefully")
	}

	logger.Info("Server stopped")
}
