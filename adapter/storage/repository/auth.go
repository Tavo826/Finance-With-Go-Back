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

func (ar *AuthRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {

	result, err := ar.db.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return user, nil
}

func (ar *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {

	var user domain.User

	if err := ar.db.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
