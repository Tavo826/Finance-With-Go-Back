package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// MongoTransactionManager implements port.TransactionManager using
// MongoDB client sessions, so multi-collection writes commit or abort together.
type MongoTransactionManager struct {
	client *mongo.Client
}

func NewMongoTransactionManager(client *mongo.Client) *MongoTransactionManager {
	return &MongoTransactionManager{client}
}

func (tm *MongoTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {

	session, err := tm.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	})

	return err
}
