package db

import (
	"context"
	"log/slog"
	"personal-finance/adapter/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(ctx context.Context, config *config.DB) (*mongo.Database, error) {

	clientOptions := options.Client().ApplyURI(config.Connection)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		slog.Error("error connecting to database", "error", err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		slog.Error("error in database communication", "error", err)
		return nil, err
	}

	return client.Database(config.Database), nil
}
