package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"personal-finance/adapter/config"
	"personal-finance/adapter/handler/http"
	"personal-finance/adapter/storage/cloud"
	"personal-finance/adapter/storage/cloud/adapter"
	"personal-finance/adapter/storage/db"
	"personal-finance/adapter/storage/db/repository"
	"personal-finance/adapter/web/mail"
	"personal-finance/core/service"

	"github.com/go-playground/validator/v10"
)

func main() {

	config, err := config.New()
	if err != nil {
		slog.Error("error loading environment variables", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting the application", "app", config.App.Name, "env", config.App.Env)

	ctx := context.Background()

	dbClient, database, err := db.New(ctx, config.DB)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}
	txManager := db.NewMongoTransactionManager(dbClient)

	storage, err := cloud.New(ctx, config.ImageCloud)
	if err != nil {
		slog.Error("Error connecting to storage service", "error", err)
		os.Exit(1)
	}

	validate := validator.New()

	originRepo := repository.NewOriginRepository(database, config.DB)
	originService := service.NewOriginService(originRepo)
	originHandler := http.NewOriginHandler(originService, validate)

	transactionRepo := repository.NewTransactionRepository(database, config.DB)
	transactionService := service.NewTransactionService(transactionRepo, originRepo, txManager)
	transactionHandler := http.NewTransactionHandler(transactionService, validate)

	authRepo := repository.NewAuthRepository(database, config.DB)
	imageAdapter := adapter.NewImageAdapter(storage)
	authService := service.NewAuthService(authRepo, transactionRepo, imageAdapter)
	authHandler := http.NewAuthHandler(authService, validate, config.Token)

	mailAdapter := mail.NewMailReportAdapter(config.Mail)
	reportService := service.NewReportService(authService, transactionService, originService, mailAdapter)
	reportHandler := http.NewReportHandler(reportService)

	router, err := http.NewRouter(config, *transactionHandler, *authHandler, *originHandler, *reportHandler)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	listenPort := fmt.Sprintf(":%s", config.App.Port)
	slog.Info("Starting HTTP sever", "port", listenPort)

	err = router.Serve(listenPort)
	if err != nil {
		slog.Error("Error starting HTTP server", "error", err)
		os.Exit(1)
	}
}
