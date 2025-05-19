package repository

import (
	"context"
	"personal-finance/adapter/config"
	"personal-finance/core/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionRepository struct {
	db *mongo.Collection
}

func NewTransactionRepository(db *mongo.Database, config *config.DB) *TransactionRepository {
	return &TransactionRepository{
		db.Collection(config.Transactions),
	}
}

func (tr *TransactionRepository) GetTransactionsByUserId(ctx context.Context, page, limit uint64, userId string) ([]domain.Transaction, any, any, error) {

	var transactions []domain.Transaction

	filter := bson.M{
		"user_id": userId,
	}
	total, err := tr.db.CountDocuments(ctx, filter)

	if err != nil {
		return nil, nil, nil, err
	}

	offset := int64((page - 1) * limit)

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetSkip(offset)
	findOptions.SetLimit(int64(limit))

	cursor, err := tr.db.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, nil, nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var transaction domain.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, nil, nil, err
		}
		transactions = append(transactions, transaction)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return transactions, total, totalPages, nil
}

func (tr *TransactionRepository) GetTransactionsByDate(
	ctx context.Context,
	userId string,
	page, limit uint64,
	year int,
	month int,
) ([]domain.Transaction, any, any, error) {

	var transactions []domain.Transaction
	var filter bson.M

	if month == 0 {

		startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

		filter = bson.M{
			"user_id": userId,
			"created_at": bson.M{
				"$gte": startDate,
				"$lt":  endDate,
			},
		}
	} else {

		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0)

		filter = bson.M{
			"user_id": userId,
			"created_at": bson.M{
				"$gte": startDate,
				"$lt":  endDate,
			},
		}
	}

	total, err := tr.db.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, nil, err
	}

	offset := int64((page - 1) * limit)

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetSkip(offset)
	findOptions.SetLimit(int64(limit))

	cursor, err := tr.db.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, nil, nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var transaction domain.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, nil, nil, err
		}
		transactions = append(transactions, transaction)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return transactions, total, totalPages, nil
}

func (tr *TransactionRepository) GetTransaction(ctx context.Context, id string) (*domain.Transaction, error) {

	var transaction domain.Transaction
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	if err := tr.db.FindOne(ctx, bson.M{"_id": objectId}).Decode(&transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (tr *TransactionRepository) CreateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {

	result, err := tr.db.InsertOne(ctx, transaction)

	if err != nil {
		return nil, err
	}

	transaction.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return transaction, nil
}

func (tr *TransactionRepository) UpdateTransaction(ctx context.Context, id string, updatedTransaction *domain.Transaction) (*domain.Transaction, error) {

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	update := bson.M{"$set": updatedTransaction}

	result, err := tr.db.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, domain.ErrDataNotFound
	}

	updatedTransaction.ID = id

	return updatedTransaction, nil
}

func (tr *TransactionRepository) DeleteTransaction(ctx context.Context, id string) error {

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	result, err := tr.db.DeleteOne(ctx, bson.M{"_id": objectId})

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrDataNotFound
	}

	return nil
}
