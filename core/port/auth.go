package port

import (
	"context"
	"personal-finance/core/domain"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, createUser *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type AuthService interface {
	CreateUser(ctx context.Context, createUser *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	VerifyUserEmail(ctx context.Context, email string) (bool, error)
}
