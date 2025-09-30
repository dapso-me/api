package main

import (
	"api/internal/configuration"
	"api/internal/transport/http"
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

	_, err = postgresql.New(&cfg.PostgreSQL)
	if err != nil {
		logger.Fatal("failed to connect db", zap.Error(err))
	}

	httpServer := http.New()
	httpServer.SetupRouter()

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
