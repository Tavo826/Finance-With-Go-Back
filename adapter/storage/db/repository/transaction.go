package repository

import (
	"context"
	"personal-finance/adapter/config"
	"personal-finance/core/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository struct {
	db *mongo.Collection
}

func NewTransactionRepository(db *mongo.Database, config *config.DB) *TransactionRepository {
	return &TransactionRepository{
		db.Collection(config.Transactions),
	}
}

func (tr *TransactionRepository) GetTransactionsByUserId(
	ctx context.Context,
	page, limit uint64,
	userId string,
) ([]domain.Transaction, any, any, error) {

	pipeline := mongo.Pipeline{

		// Filter by usuario
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},

		// origin_id to ObjectId
		{{Key: "$addFields", Value: bson.D{
			{Key: "origin_object_id", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$and", Value: bson.A{
						bson.D{{Key: "$ifNull", Value: bson.A{"$origin_id", false}}},
						bson.D{{Key: "$ne", Value: bson.A{"$origin_id", ""}}},
					}}},
					bson.D{{Key: "$toObjectId", Value: "$origin_id"}},
					nil,
				}},
			}},
		}}},

		// Left join origins
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "origins"},
			{Key: "localField", Value: "origin_object_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "origin"},
		}}},

		// Array origin to object
		{{Key: "$addFields", Value: bson.D{
			{Key: "origin", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$gt", Value: bson.A{bson.D{{Key: "$size", Value: "$origin"}}, 0}}},
					bson.D{{Key: "$arrayElemAt", Value: bson.A{"$origin", 0}}},
					nil,
				}},
			}},
		}}},

		// Clear temporal field
		{{Key: "$unset", Value: "origin_object_id"}},

		// Order by creation date DESC
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},

		// Pagination
		{{Key: "$skip", Value: int64((page - 1) * limit)}},
		{{Key: "$limit", Value: int64(limit)}},
	}

	return tr.findtransactionUsingPipeline(ctx, pipeline, userId, limit)
}

func (tr *TransactionRepository) GetTransactionsByDate(
	ctx context.Context,
	userId string,
	page, limit uint64,
	year int,
	month int,
) ([]domain.Transaction, any, any, error) {

	var dateFilter bson.M

	if month == 0 {

		startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

		dateFilter = bson.M{
			"user_id": userId,
			"created_at": bson.M{
				"$gte": startDate,
				"$lt":  endDate,
			},
		}
	} else {

		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0)

		dateFilter = bson.M{
			"user_id": userId,
			"created_at": bson.M{
				"$gte": startDate,
				"$lt":  endDate,
			},
		}
	}

	pipeline := mongo.Pipeline{

		// Filter by user and date
		{{Key: "$match", Value: dateFilter}},

		// origin_id to ObjectId
		{{Key: "$addFields", Value: bson.D{
			{Key: "origin_object_id", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$and", Value: bson.A{
						bson.D{{Key: "$ifNull", Value: bson.A{"$origin_id", false}}},
						bson.D{{Key: "$ne", Value: bson.A{"$origin_id", ""}}},
					}}},
					bson.D{{Key: "$toObjectId", Value: "$origin_id"}},
					nil,
				}},
			}},
		}}},

		// Left join origins
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "origins"},
			{Key: "localField", Value: "origin_object_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "origin"},
		}}},

		// Array origin to object
		{{Key: "$addFields", Value: bson.D{
			{Key: "origin", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$gt", Value: bson.A{bson.D{{Key: "$size", Value: "$origin"}}, 0}}},
					bson.D{{Key: "$arrayElemAt", Value: bson.A{"$origin", 0}}},
					nil,
				}},
			}},
		}}},

		// Clear temporal field
		{{Key: "$unset", Value: "origin_object_id"}},

		// Order by creation date DESC
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},

		// Pagination
		{{Key: "$skip", Value: int64((page - 1) * limit)}},
		{{Key: "$limit", Value: int64(limit)}},
	}

	return tr.findtransactionUsingPipeline(ctx, pipeline, userId, limit)
}

func (tr *TransactionRepository) GetTransactionsBySubject(
	ctx context.Context,
	userId string,
	page, limit uint64,
	subject string,
	personOrBusiness string,
) ([]domain.Transaction, any, any, error) {

	var subjectFilter bson.M

	if personOrBusiness == "" {

		subjectFilter = bson.M{
			"user_id": userId,
			"subject": subject,
		}
	} else {

		subjectFilter = bson.M{
			"user_id":         userId,
			"subject":         subject,
			"person_business": personOrBusiness,
		}
	}

	pipeline := mongo.Pipeline{

		// Filter by user and date
		{{Key: "$match", Value: subjectFilter}},

		// origin_id to ObjectId
		{{Key: "$addFields", Value: bson.D{
			{Key: "origin_object_id", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$and", Value: bson.A{
						bson.D{{Key: "$ifNull", Value: bson.A{"$origin_id", false}}},
						bson.D{{Key: "$ne", Value: bson.A{"$origin_id", ""}}},
					}}},
					bson.D{{Key: "$toObjectId", Value: "$origin_id"}},
					nil,
				}},
			}},
		}}},

		// Left join origins
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "origins"},
			{Key: "localField", Value: "origin_object_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "origin"},
		}}},

		// Array origin to object
		{{Key: "$addFields", Value: bson.D{
			{Key: "origin", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$gt", Value: bson.A{bson.D{{Key: "$size", Value: "$origin"}}, 0}}},
					bson.D{{Key: "$arrayElemAt", Value: bson.A{"$origin", 0}}},
					nil,
				}},
			}},
		}}},

		// Clear temporal field
		{{Key: "$unset", Value: "origin_object_id"}},

		// Order by creation date DESC
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},

		// Pagination
		{{Key: "$skip", Value: int64((page - 1) * limit)}},
		{{Key: "$limit", Value: int64(limit)}},
	}

	return tr.findtransactionUsingPipeline(ctx, pipeline, userId, limit)
}

func (tr *TransactionRepository) GetTransactionById(ctx context.Context, id string) (*domain.Transaction, error) {

	var transaction domain.Transaction
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{

		// Filter by id
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectId}}}},

		// origin_id to ObjectId
		{{Key: "$addFields", Value: bson.D{
			{Key: "origin_object_id", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$and", Value: bson.A{
						bson.D{{Key: "$ifNull", Value: bson.A{"$origin_id", false}}},
						bson.D{{Key: "$ne", Value: bson.A{"$origin_id", ""}}},
					}}},
					bson.D{{Key: "$toObjectId", Value: "$origin_id"}},
					nil,
				}},
			}},
		}}},

		// Left join origins
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "origins"},
			{Key: "localField", Value: "origin_object_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "origin"},
		}}},

		// Array origin to object
		{{Key: "$addFields", Value: bson.D{
			{Key: "origin", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$gt", Value: bson.A{bson.D{{Key: "$size", Value: "$origin"}}, 0}}},
					bson.D{{Key: "$arrayElemAt", Value: bson.A{"$origin", 0}}},
					nil,
				}},
			}},
		}}},

		// Clear temporal field
		{{Key: "$unset", Value: "origin_object_id"}},
	}

	cursor, err := tr.db.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		return &transaction, nil
	}

	return nil, domain.ErrDataNotFound
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

func (tr *TransactionRepository) DeleteTransactionsByUserId(ctx context.Context, id string) error {

	_, err := tr.db.DeleteMany(ctx, bson.M{"user_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (tr *TransactionRepository) findtransactionUsingPipeline(
	ctx context.Context,
	pipeline mongo.Pipeline,
	userId string,
	limit uint64,
) ([]domain.Transaction, any, any, error) {

	var transactions []domain.Transaction

	total, err := tr.db.CountDocuments(ctx, bson.M{"user_id": userId})
	if err != nil {
		return nil, nil, nil, err
	}

	cursor, err := tr.db.Aggregate(ctx, pipeline)
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

	if err := cursor.Err(); err != nil {
		return nil, nil, nil, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return transactions, total, totalPages, nil
}
