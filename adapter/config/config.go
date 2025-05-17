package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App   *App
		DB    *DB
		Token *Token
	}

	App struct {
		Name           string
		Env            string
		Port           string
		AllowedOrigins string
	}

	DB struct {
		Connection   string
		Database     string
		Transactions string
		Users        string
	}

	Token struct {
		JwtSecret string
	}
)

func New() (*Container, error) {

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	app := &App{
		Name:           os.Getenv("APP_NAME"),
		Env:            os.Getenv("APP_ENV"),
		Port:           os.Getenv("PORT"),
		AllowedOrigins: os.Getenv("ALLOWED_ORIGINS"),
	}

	db := &DB{
		Connection:   os.Getenv("MONGO_CONNECTION_STRING"),
		Database:     os.Getenv("MONGO_DATABASE_NAME"),
		Transactions: os.Getenv("MONGO_COLLECTION_TRANSACTION"),
		Users:        os.Getenv("MONGO_COLLECTION_USER"),
	}

	token := &Token{
		JwtSecret: os.Getenv("JWT_SECRET"),
	}

	return &Container{
		app,
		db,
		token,
	}, nil
}
