package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App        *App
		DB         *DB
		ImageCloud *ImageCloud
		Mail       *Mail
		Token      *Token
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
		Origin       string
		Users        string
	}

	ImageCloud struct {
		Connection string
	}

	Mail struct {
		Host     string
		Username string
		Password string
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
		Origin:       os.Getenv("MONGO_COLLECTION_ORIGIN"),
		Users:        os.Getenv("MONGO_COLLECTION_USER"),
	}

	imageCloud := &ImageCloud{
		Connection: os.Getenv("CLOUDINARY_URL"),
	}

	mailService := &Mail{
		Host:     os.Getenv("MAIL_SERVICE_HOST"),
		Username: os.Getenv("MAIL_SERVICE_USERNAME"),
		Password: os.Getenv("MAIL_SERVICE_PASS"),
	}

	token := &Token{
		JwtSecret: os.Getenv("JWT_SECRET"),
	}

	return &Container{
		app,
		db,
		imageCloud,
		mailService,
		token,
	}, nil
}
