package port

import (
	"context"
	"mime/multipart"
	"personal-finance/core/domain"
)

type AuthRepository interface {
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, createUser *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, updateUser *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type AuthService interface {
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	VerifyUserEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, createUser *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, updateUser *domain.User) (*domain.User, error)
	UpdateUserProfileImage(ctx context.Context, file multipart.File, userId string) (*domain.Image, error)
	DeleteUserProfileImage(ctx context.Context, publicId string) error
	DeleteUser(ctx context.Context, id string) error
	DeleteTransactionsByUserId(ctx context.Context, id string) error
}
