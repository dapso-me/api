package main

import (
	customer_usecase "api/internal/application/customer"
	"api/internal/configuration"
	customer_postgresql_repository "api/internal/infrastructure/repository/customer/postgresql"
	session_postgresql_repository "api/internal/infrastructure/repository/session/postgresql"
	"api/internal/transport/http"
	account_handler "api/internal/transport/http/handlers/account"
	"api/pkg/logger"
	"api/pkg/postgresql"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	cfg, err := configuration.Load()
	if err != nil {
		log.Fatal(err)
	}

	loggerConfig := logger.Config{
		OutputPath:      cfg.Logger.OutputPath,
		ErrorOutputPath: cfg.Logger.ErrorOutputPath,
	}
	logger, err := logger.New(loggerConfig)
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgresql.New(&cfg.PostgreSQL)
	if err != nil {
		logger.Fatal("failed to connect db", zap.Error(err))
	}

	// repositories
	customerRepoPG := customer_postgresql_repository.New(db)
	sessionRepoPG := session_postgresql_repository.New(db)

	// usecases
	customerUC := customer_usecase.New(customerRepoPG, sessionRepoPG)

	// handlers
	accountHandler := account_handler.New(logger, customerUC)

	httpServer := http.New(sessionRepoPG)
	httpServer.SetupRouter(
		accountHandler,
	)

	go func() {
		if err := httpServer.Run("0.0.0.0:4000"); err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}
}
