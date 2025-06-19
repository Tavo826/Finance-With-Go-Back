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

type OriginRepository struct {
	db *mongo.Collection
}

func NewOriginRepository(db *mongo.Database, config *config.DB) *OriginRepository {
	return &OriginRepository{
		db.Collection(config.Origin),
	}
}

func (or *OriginRepository) GetOriginsByUserId(ctx context.Context, userId string) ([]domain.Origin, error) {

	var origins []domain.Origin

	filter := bson.M{
		"user_id": userId,
	}

	cursor, err := or.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var origin domain.Origin
		if err := cursor.Decode(&origin); err != nil {
			return nil, err
		}
		origins = append(origins, origin)
	}

	return origins, nil
}

func (or *OriginRepository) GetOriginById(ctx context.Context, id string) (*domain.Origin, error) {

	var origin domain.Origin
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	if err := or.db.FindOne(ctx, bson.M{"_id": objectId}).Decode(&origin); err != nil {
		return nil, err
	}

	return &origin, nil
}

func (or *OriginRepository) CreateOrigin(ctx context.Context, origin *domain.Origin) (*domain.Origin, error) {

	result, err := or.db.InsertOne(ctx, origin)
	if err != nil {
		return nil, err
	}

	origin.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return origin, nil
}

func (or *OriginRepository) UpdateOrigin(ctx context.Context, id string, updatedOrigin *domain.Origin) (*domain.Origin, error) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	origin := domain.OriginRequest{
		UserId:    updatedOrigin.UserId,
		Name:      updatedOrigin.Name,
		Total:     updatedOrigin.Total,
		CreatedAt: updatedOrigin.CreatedAt,
		UpdatedAt: time.Now(),
	}

	update := bson.M{"$set": origin}

	result, err := or.db.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, domain.ErrDataNotFound
	}

	updatedOrigin.ID = id

	return updatedOrigin, nil
}

func (or *OriginRepository) DeleteOrigin(ctx context.Context, id string) error {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := or.db.DeleteOne(ctx, bson.M{"_id": objectId})

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrDataNotFound
	}

	return nil
}
