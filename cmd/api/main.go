package main

import (
	"api/internal/configuration"
	"api/internal/transport/http"
	"api/pkg/logger"
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

	// dbConfig := postgresql.Config{
	// 	Host:     cfg.PostgreSQL.Host,
	// 	Port:     cfg.PostgreSQL.Port,
	// 	User:     cfg.PostgreSQL.User,
	// 	Password: cfg.PostgreSQL.Pass,
	// 	Name:     cfg.PostgreSQL.Name,
	// 	SSLMode:  cfg.PostgreSQL.SSLMode,
	// }
	// db, err := postgresql.New(&dbConfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// db.Ping()

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
