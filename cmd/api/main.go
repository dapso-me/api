package main

import (
	category_usecase "api/internal/application/category"
	customer_usecase "api/internal/application/customer"
	dish_usecase "api/internal/application/dish"
	project_usecase "api/internal/application/project"
	"api/internal/configuration"
	category_postgresql_repository "api/internal/infrastructure/repository/category/postgresql"
	customer_postgresql_repository "api/internal/infrastructure/repository/customer/postgresql"
	dish_postgresql_repository "api/internal/infrastructure/repository/dish/postgresql"
	project_postgresql_repository "api/internal/infrastructure/repository/project/postgresql"
	session_postgresql_repository "api/internal/infrastructure/repository/session/postgresql"
	"api/internal/transport/http"
	account_handler "api/internal/transport/http/handlers/account"
	menu_handler "api/internal/transport/http/handlers/menu"
	project_handler "api/internal/transport/http/handlers/project"
	"api/pkg/logger"
	"api/pkg/postgresql"
	"api/pkg/storage"
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

	storage, err := storage.New(&cfg.S3)
	if err != nil {
		logger.Fatal("failed to connect db", zap.Error(err))
	}

	// repositories
	customerRepoPG := customer_postgresql_repository.New(db)
	sessionRepoPG := session_postgresql_repository.New(db)
	projectRepoPG := project_postgresql_repository.New(db)
	categoryRepoPG := category_postgresql_repository.New(db)
	dishRepoPG := dish_postgresql_repository.New(db)

	// usecases
	customerUC := customer_usecase.New(customerRepoPG, sessionRepoPG)
	projectUC := project_usecase.New(customerUC, projectRepoPG, categoryRepoPG)
	categoryUC := category_usecase.New(categoryRepoPG, projectRepoPG, dishRepoPG)
	dishUC := dish_usecase.New(dishRepoPG, projectUC, storage)

	// handlers
	accountHandler := account_handler.New(logger, customerUC)
	projectHandler := project_handler.New(logger, projectUC)
	menuHandler := menu_handler.New(logger, categoryUC, dishUC)

	httpServer := http.New(sessionRepoPG)
	httpServer.SetupRouter(
		accountHandler,
		projectHandler,
		menuHandler,
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
