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
	"personal-finance/core/service"

	"github.com/gin-gonic/gin"
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

	db, err := db.New(ctx, config.DB)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}

	storage, err := cloud.New(ctx, config.ImageCloud)
	if err != nil {
		slog.Error("Error connecting to storage service", "error", err)
		os.Exit(1)
	}

	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	validate := validator.New()

	transactionRepo := repository.NewTransactionRepository(db, config.DB)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionHandler := http.NewTransactionHandler(transactionService, validate)

	authRepo := repository.NewAuthRepository(db, config.DB)
	imageAdapter := adapter.NewImageAdapter(storage)
	authService := service.NewAuthService(authRepo, transactionRepo, imageAdapter)
	authHandler := http.NewAuthHandler(authService, validate, config.Token)

	router, err := http.NewRouter(config, *transactionHandler, *authHandler)
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
