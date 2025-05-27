package port

import (
	"context"
	"personal-finance/core/domain"
)

type AuthRepository interface {
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, createUser *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, updateUser *domain.User) (*domain.User, error)
}

type AuthService interface {
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	VerifyUserEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, createUser *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, updateUser *domain.User) (*domain.User, error)
}
