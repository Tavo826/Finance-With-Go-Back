package service

import (
	"context"
	"log"
	"personal-finance/core/domain"
	"personal-finance/core/port"
)

type AuthService struct {
	repo port.AuthRepository
}

func NewAuthService(repo port.AuthRepository) *AuthService {
	return &AuthService{
		repo,
	}
}

func (as *AuthService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {

	user, err := as.repo.CreateUser(ctx, user)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (as *AuthService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {

	user, err := as.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

func (as *AuthService) VerifyUserEmail(ctx context.Context, email string) (bool, error) {

	_, err := as.repo.GetUserByEmail(ctx, email)

	if err != nil {
		if err == domain.ErrDataNotFound {
			return false, nil
		}
		if err.Error() == domain.ErrNoDocuments.Error() {
			return false, nil
		}
		log.Println(err.Error())
		return true, domain.ErrInternal
	}

	return true, nil
}
