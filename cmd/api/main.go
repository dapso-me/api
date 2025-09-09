package main

import (
	customer_usecase "api/internal/application/customer"
	otp_usecase "api/internal/application/otp"
	session_usecase "api/internal/application/session"
	"api/internal/configuration"
	customer_repository "api/internal/infrastructure/repository/postgresql/customer"
	otp_postgresql_repository "api/internal/infrastructure/repository/postgresql/otp"
	session_postgresql_repository "api/internal/infrastructure/repository/postgresql/session"
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

	customerRepo := customer_repository.New(db)
	sessionRepo := session_postgresql_repository.New(db)
	otpRepo := otp_postgresql_repository.New(db)

	otpUseCase := otp_usecase.New(otpRepo, nil)
	sessionUseCase := session_usecase.New(sessionRepo)
	customerUseCase := customer_usecase.New(customerRepo, otpUseCase, sessionUseCase)

	accountHandler := account_handler.New(customerUseCase)

	httpServer := http.New()
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
