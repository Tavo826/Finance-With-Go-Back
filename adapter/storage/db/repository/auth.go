package repository

import (
	"context"
	"personal-finance/adapter/config"
	"personal-finance/core/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	db *mongo.Collection
}

func NewAuthRepository(db *mongo.Database, config *config.DB) *AuthRepository {
	return &AuthRepository{
		db.Collection(config.Users),
	}
}

func (ar *AuthRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {

	var users []domain.User

	cursor, err := ar.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (ar *AuthRepository) GetUserById(ctx context.Context, id string) (*domain.User, error) {

	var user domain.User
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err := ar.db.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ar *AuthRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {

	result, err := ar.db.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return user, nil
}

func (ar *AuthRepository) UpdateUser(ctx context.Context, id string, updatedUser *domain.User) (*domain.User, error) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{"$set": updatedUser}

	result, err := ar.db.UpdateOne(ctx, bson.M{"_id": objectId}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, domain.ErrDataNotFound
	}

	updatedUser.ID = id

	return updatedUser, nil
}

func (ar *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {

	var user domain.User

	if err := ar.db.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ar *AuthRepository) DeleteUser(ctx context.Context, id string) error {

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	result, err := ar.db.DeleteOne(ctx, bson.M{"_id": objectId})

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrDataNotFound
	}

	return nil
}
